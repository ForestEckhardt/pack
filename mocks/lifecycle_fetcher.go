// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpack/pack (interfaces: LifecycleFetcher)

// Package mocks is a generated GoMock package.
package mocks

import (
	semver "github.com/Masterminds/semver"
	lifecycle "github.com/buildpack/pack/lifecycle"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLifecycleFetcher is a mock of LifecycleFetcher interface
type MockLifecycleFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockLifecycleFetcherMockRecorder
}

// MockLifecycleFetcherMockRecorder is the mock recorder for MockLifecycleFetcher
type MockLifecycleFetcherMockRecorder struct {
	mock *MockLifecycleFetcher
}

// NewMockLifecycleFetcher creates a new mock instance
func NewMockLifecycleFetcher(ctrl *gomock.Controller) *MockLifecycleFetcher {
	mock := &MockLifecycleFetcher{ctrl: ctrl}
	mock.recorder = &MockLifecycleFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLifecycleFetcher) EXPECT() *MockLifecycleFetcherMockRecorder {
	return m.recorder
}

// Fetch mocks base method
func (m *MockLifecycleFetcher) Fetch(arg0 *semver.Version, arg1 string) (lifecycle.Metadata, error) {
	ret := m.ctrl.Call(m, "Fetch", arg0, arg1)
	ret0, _ := ret[0].(lifecycle.Metadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch
func (mr *MockLifecycleFetcherMockRecorder) Fetch(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockLifecycleFetcher)(nil).Fetch), arg0, arg1)
}
