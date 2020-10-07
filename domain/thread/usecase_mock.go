// Code generated by MockGen. DO NOT EDIT.
// Source: ./domain/thread/usecase.go

// Package thread is a generated GoMock package.
package thread

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// get mocks base method
func (m *MockRepository) get(arg0 context.Context, arg1 repositoryGetRequest) ([]Thread, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "get", arg0, arg1)
	ret0, _ := ret[0].([]Thread)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// get indicates an expected call of get
func (mr *MockRepositoryMockRecorder) get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "get", reflect.TypeOf((*MockRepository)(nil).get), arg0, arg1)
}

// create mocks base method
func (m *MockRepository) create(arg0 context.Context, arg1 repositoryCreateRequest) (Thread, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "create", arg0, arg1)
	ret0, _ := ret[0].(Thread)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// create indicates an expected call of create
func (mr *MockRepositoryMockRecorder) create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "create", reflect.TypeOf((*MockRepository)(nil).create), arg0, arg1)
}

// update mocks base method
func (m *MockRepository) update(arg0 context.Context, arg1 repositoryUpdateRequest) (Thread, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "update", arg0, arg1)
	ret0, _ := ret[0].(Thread)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// update indicates an expected call of update
func (mr *MockRepositoryMockRecorder) update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "update", reflect.TypeOf((*MockRepository)(nil).update), arg0, arg1)
}

// open mocks base method
func (m *MockRepository) open(arg0 context.Context, arg1 repositoryOpenRequest) (Thread, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "open", arg0, arg1)
	ret0, _ := ret[0].(Thread)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// open indicates an expected call of open
func (mr *MockRepositoryMockRecorder) open(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "open", reflect.TypeOf((*MockRepository)(nil).open), arg0, arg1)
}

// close mocks base method
func (m *MockRepository) close(arg0 context.Context, arg1 repositoryCloseRequest) (Thread, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "close", arg0, arg1)
	ret0, _ := ret[0].(Thread)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// close indicates an expected call of close
func (mr *MockRepositoryMockRecorder) close(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "close", reflect.TypeOf((*MockRepository)(nil).close), arg0, arg1)
}
