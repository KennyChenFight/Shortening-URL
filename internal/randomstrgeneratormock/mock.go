// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KennyChenFight/randstr (interfaces: RandomStrGenerator)

// Package randomstrgeneratormock is a generated GoMock package.
package randomstrgeneratormock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRandomStrGenerator is a mock of RandomStrGenerator interface.
type MockRandomStrGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockRandomStrGeneratorMockRecorder
}

// MockRandomStrGeneratorMockRecorder is the mock recorder for MockRandomStrGenerator.
type MockRandomStrGeneratorMockRecorder struct {
	mock *MockRandomStrGenerator
}

// NewMockRandomStrGenerator creates a new mock instance.
func NewMockRandomStrGenerator(ctrl *gomock.Controller) *MockRandomStrGenerator {
	mock := &MockRandomStrGenerator{ctrl: ctrl}
	mock.recorder = &MockRandomStrGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRandomStrGenerator) EXPECT() *MockRandomStrGeneratorMockRecorder {
	return m.recorder
}

// GenerateRandomStr mocks base method.
func (m *MockRandomStrGenerator) GenerateRandomStr(arg0 int) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRandomStr", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GenerateRandomStr indicates an expected call of GenerateRandomStr.
func (mr *MockRandomStrGeneratorMockRecorder) GenerateRandomStr(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRandomStr", reflect.TypeOf((*MockRandomStrGenerator)(nil).GenerateRandomStr), arg0)
}
