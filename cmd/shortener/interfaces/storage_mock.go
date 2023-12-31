// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nikishin42/shortener/cmd/shortener/interfaces (interfaces: Storage)

// Package interfaces is a generated GoMock package.
package interfaces

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetFullURL mocks base method.
func (m *MockStorage) GetFullURL(arg0 string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFullURL", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetFullURL indicates an expected call of GetFullURL.
func (mr *MockStorageMockRecorder) GetFullURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFullURL", reflect.TypeOf((*MockStorage)(nil).GetFullURL), arg0)
}

// GetID mocks base method.
func (m *MockStorage) GetID(arg0 string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetID indicates an expected call of GetID.
func (mr *MockStorageMockRecorder) GetID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetID", reflect.TypeOf((*MockStorage)(nil).GetID), arg0)
}

// SetPair mocks base method.
func (m *MockStorage) SetPair(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPair", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPair indicates an expected call of SetPair.
func (mr *MockStorageMockRecorder) SetPair(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPair", reflect.TypeOf((*MockStorage)(nil).SetPair), arg0, arg1)
}
