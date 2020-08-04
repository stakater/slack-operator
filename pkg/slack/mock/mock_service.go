// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stakater/slack-operator/pkg/slack (interfaces: Service)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateChannel mocks base method
func (m *MockService) CreateChannel(arg0 string, arg1 bool) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChannel", arg0, arg1)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChannel indicates an expected call of CreateChannel
func (mr *MockServiceMockRecorder) CreateChannel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChannel", reflect.TypeOf((*MockService)(nil).CreateChannel), arg0, arg1)
}

// InviteUsers mocks base method
func (m *MockService) InviteUsers(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InviteUsers", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InviteUsers indicates an expected call of InviteUsers
func (mr *MockServiceMockRecorder) InviteUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InviteUsers", reflect.TypeOf((*MockService)(nil).InviteUsers), arg0, arg1)
}

// RenameChannel mocks base method
func (m *MockService) RenameChannel(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenameChannel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenameChannel indicates an expected call of RenameChannel
func (mr *MockServiceMockRecorder) RenameChannel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenameChannel", reflect.TypeOf((*MockService)(nil).RenameChannel), arg0, arg1)
}

// SetDescription mocks base method
func (m *MockService) SetDescription(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDescription", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDescription indicates an expected call of SetDescription
func (mr *MockServiceMockRecorder) SetDescription(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDescription", reflect.TypeOf((*MockService)(nil).SetDescription), arg0, arg1)
}

// SetTopic mocks base method
func (m *MockService) SetTopic(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTopic", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTopic indicates an expected call of SetTopic
func (mr *MockServiceMockRecorder) SetTopic(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTopic", reflect.TypeOf((*MockService)(nil).SetTopic), arg0, arg1)
}
