package pack_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/buildpack/imgutil/fakes"
	"github.com/golang/mock/gomock"
	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpack/pack"
	"github.com/buildpack/pack/api"
	"github.com/buildpack/pack/buildpackage"
	"github.com/buildpack/pack/dist"
	ifakes "github.com/buildpack/pack/internal/fakes"
	"github.com/buildpack/pack/internal/logging"
	h "github.com/buildpack/pack/testhelpers"
	"github.com/buildpack/pack/testmocks"
)

func TestCreateBuildpackage(t *testing.T) {
	color.Disable(true)
	defer func() { color.Disable(false) }()
	spec.Run(t, "CreateBuildpackage", testCreateBuildpackage, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testCreateBuildpackage(t *testing.T, when spec.G, it spec.S) {
	var (
		client                *pack.Client
		mockController        *gomock.Controller
		mockDownloader        *testmocks.MockDownloader
		mockImageFactory      *testmocks.MockImageFactory
		fakeBuildpackageImage *fakes.Image
		out                   bytes.Buffer
	)

	it.Before(func() {
		mockController = gomock.NewController(t)
		mockDownloader = testmocks.NewMockDownloader(mockController)
		mockImageFactory = testmocks.NewMockImageFactory(mockController)

		fakeBuildpackageImage = fakes.NewImage("some/package", "", "")
		mockImageFactory.EXPECT().NewImage("some/package", true).Return(fakeBuildpackageImage, nil).AnyTimes()

		var err error
		client, err = pack.NewClient(
			pack.WithLogger(logging.NewLogWithWriters(&out, &out)),
			pack.WithDownloader(mockDownloader),
			pack.WithImageFactory(mockImageFactory),
		)
		h.AssertNil(t, err)
	})

	it.After(func() {
		fakeBuildpackageImage.Cleanup()
		mockController.Finish()
	})

	when("#CreateBuildpackage", func() {
		when("buildpackage config is valid", func() {
			var opts pack.CreateBuildpackageOptions

			it.Before(func() {
				opts = pack.CreateBuildpackageOptions{
					Name: fakeBuildpackageImage.Name(),
					Config: buildpackage.Config{
						Default: dist.BuildpackInfo{
							ID:      "bp.one",
							Version: "1.2.3",
						},
						Blobs: []dist.BlobConfig{
							{URI: "https://example.com/bp.one.tgz"},
						},
						Stacks: []dist.Stack{
							{ID: "some.stack.id"},
						},
					},
				}

				buildpack, err := ifakes.NewBuildpackFromDescriptor(dist.BuildpackDescriptor{
					API: api.MustParse("0.2"),
					Info: dist.BuildpackInfo{
						ID:      "bp.one",
						Version: "1.2.3",
					},
					Stacks: []dist.Stack{
						{ID: "some.stack.id"},
					},
					Order: nil,
				}, 0644)

				h.AssertNil(t, err)

				mockDownloader.EXPECT().Download(gomock.Any(), "https://example.com/bp.one.tgz").Return(buildpack, nil).AnyTimes()
			})

			it("sets metadata", func() {
				h.AssertNil(t, client.CreateBuildpackage(context.TODO(), opts))
				h.AssertEq(t, fakeBuildpackageImage.IsSaved(), true)

				labelData, err := fakeBuildpackageImage.Label("io.buildpacks.buildpackage.metadata")
				h.AssertNil(t, err)
				var md buildpackage.Metadata
				h.AssertNil(t, json.Unmarshal([]byte(labelData), &md))

				h.AssertEq(t, md.ID, "bp.one")
				h.AssertEq(t, md.Version, "1.2.3")
				h.AssertEq(t, len(md.Stacks), 1)
				h.AssertEq(t, md.Stacks[0].ID, "some.stack.id")
			})

			it("adds buildpack layers", func() {
				h.AssertNil(t, client.CreateBuildpackage(context.TODO(), opts))
				h.AssertEq(t, fakeBuildpackageImage.IsSaved(), true)

				dirPath := fmt.Sprintf("/cnb/buildpacks/%s/%s", "bp.one", "1.2.3")
				layerTar, err := fakeBuildpackageImage.FindLayerWithPath(dirPath)
				h.AssertNil(t, err)

				h.AssertOnTarEntry(t, layerTar, dirPath,
					h.IsDirectory(),
				)

				h.AssertOnTarEntry(t, layerTar, dirPath+"/bin/build",
					h.ContentEquals("build-contents"),
					h.HasOwnerAndGroup(0, 0),
					h.HasFileMode(0644),
				)

				h.AssertOnTarEntry(t, layerTar, dirPath+"/bin/detect",
					h.ContentEquals("detect-contents"),
					h.HasOwnerAndGroup(0, 0),
					h.HasFileMode(0644),
				)
			})

			when("when publish is true", func() {
				var fakeRemoteBuildpackageImage *fakes.Image

				it.Before(func() {
					fakeRemoteBuildpackageImage = fakes.NewImage("some/package", "", "")
					mockImageFactory.EXPECT().NewImage("some/package", false).Return(fakeRemoteBuildpackageImage, nil).AnyTimes()

					opts.Publish = true
				})

				it.After(func() {
					fakeRemoteBuildpackageImage.Cleanup()
				})

				it("saves remote image", func() {
					h.AssertNil(t, client.CreateBuildpackage(context.TODO(), opts))
					h.AssertEq(t, fakeRemoteBuildpackageImage.IsSaved(), true)
				})
			})
		})
	})
}
