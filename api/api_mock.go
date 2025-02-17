// Code generated by MockGen. DO NOT EDIT.
// Source: api/api.go

// Package api is a generated GoMock package.
package api

import (
	http "net/http"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockAPIServerInterface is a mock of APIServerInterface interface.
type MockAPIServerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAPIServerInterfaceMockRecorder
}

// MockAPIServerInterfaceMockRecorder is the mock recorder for MockAPIServerInterface.
type MockAPIServerInterfaceMockRecorder struct {
	mock *MockAPIServerInterface
}

// NewMockAPIServerInterface creates a new mock instance.
func NewMockAPIServerInterface(ctrl *gomock.Controller) *MockAPIServerInterface {
	mock := &MockAPIServerInterface{ctrl: ctrl}
	mock.recorder = &MockAPIServerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPIServerInterface) EXPECT() *MockAPIServerInterfaceMockRecorder {
	return m.recorder
}

// RegisterRoute mocks base method.
func (m *MockAPIServerInterface) RegisterRoute(path string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	m.ctrl.T.Helper()
	varargs := []interface{}{path, handler}
	for _, a := range methods {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "RegisterRoute", varargs...)
}

// RegisterRoute indicates an expected call of RegisterRoute.
func (mr *MockAPIServerInterfaceMockRecorder) RegisterRoute(path, handler interface{}, methods ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path, handler}, methods...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRoute", reflect.TypeOf((*MockAPIServerInterface)(nil).RegisterRoute), varargs...)
}

// Start mocks base method.
func (m *MockAPIServerInterface) Start(timeouts ...time.Duration) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range timeouts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Start", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockAPIServerInterfaceMockRecorder) Start(timeouts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockAPIServerInterface)(nil).Start), timeouts...)
}
