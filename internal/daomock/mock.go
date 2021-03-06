// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KennyChenFight/Shortening-URL/pkg/dao (interfaces: CacheDAO,KeyDAO,UrlDAO)

// Package daomock is a generated GoMock package.
package daomock

import (
	reflect "reflect"

	business "github.com/KennyChenFight/Shortening-URL/pkg/business"
	dao "github.com/KennyChenFight/Shortening-URL/pkg/dao"
	gomock "github.com/golang/mock/gomock"
)

// MockCacheDAO is a mock of CacheDAO interface.
type MockCacheDAO struct {
	ctrl     *gomock.Controller
	recorder *MockCacheDAOMockRecorder
}

// MockCacheDAOMockRecorder is the mock recorder for MockCacheDAO.
type MockCacheDAOMockRecorder struct {
	mock *MockCacheDAO
}

// NewMockCacheDAO creates a new mock instance.
func NewMockCacheDAO(ctrl *gomock.Controller) *MockCacheDAO {
	mock := &MockCacheDAO{ctrl: ctrl}
	mock.recorder = &MockCacheDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheDAO) EXPECT() *MockCacheDAOMockRecorder {
	return m.recorder
}

// AddOriginalURLIDInFilters mocks base method.
func (m *MockCacheDAO) AddOriginalURLIDInFilters(arg0 string) *business.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOriginalURLIDInFilters", arg0)
	ret0, _ := ret[0].(*business.Error)
	return ret0
}

// AddOriginalURLIDInFilters indicates an expected call of AddOriginalURLIDInFilters.
func (mr *MockCacheDAOMockRecorder) AddOriginalURLIDInFilters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOriginalURLIDInFilters", reflect.TypeOf((*MockCacheDAO)(nil).AddOriginalURLIDInFilters), arg0)
}

// DeleteMultiOriginalURL mocks base method.
func (m *MockCacheDAO) DeleteMultiOriginalURL(arg0 []string) *business.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMultiOriginalURL", arg0)
	ret0, _ := ret[0].(*business.Error)
	return ret0
}

// DeleteMultiOriginalURL indicates an expected call of DeleteMultiOriginalURL.
func (mr *MockCacheDAOMockRecorder) DeleteMultiOriginalURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMultiOriginalURL", reflect.TypeOf((*MockCacheDAO)(nil).DeleteMultiOriginalURL), arg0)
}

// DeleteMultiOriginalURLIDInFilters mocks base method.
func (m *MockCacheDAO) DeleteMultiOriginalURLIDInFilters(arg0 []string) (bool, *business.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMultiOriginalURLIDInFilters", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*business.Error)
	return ret0, ret1
}

// DeleteMultiOriginalURLIDInFilters indicates an expected call of DeleteMultiOriginalURLIDInFilters.
func (mr *MockCacheDAOMockRecorder) DeleteMultiOriginalURLIDInFilters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMultiOriginalURLIDInFilters", reflect.TypeOf((*MockCacheDAO)(nil).DeleteMultiOriginalURLIDInFilters), arg0)
}

// DeleteOriginalURL mocks base method.
func (m *MockCacheDAO) DeleteOriginalURL(arg0 string) *business.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOriginalURL", arg0)
	ret0, _ := ret[0].(*business.Error)
	return ret0
}

// DeleteOriginalURL indicates an expected call of DeleteOriginalURL.
func (mr *MockCacheDAOMockRecorder) DeleteOriginalURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOriginalURL", reflect.TypeOf((*MockCacheDAO)(nil).DeleteOriginalURL), arg0)
}

// DeleteOriginalURLIDInFilters mocks base method.
func (m *MockCacheDAO) DeleteOriginalURLIDInFilters(arg0 string) (bool, *business.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOriginalURLIDInFilters", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*business.Error)
	return ret0, ret1
}

// DeleteOriginalURLIDInFilters indicates an expected call of DeleteOriginalURLIDInFilters.
func (mr *MockCacheDAOMockRecorder) DeleteOriginalURLIDInFilters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOriginalURLIDInFilters", reflect.TypeOf((*MockCacheDAO)(nil).DeleteOriginalURLIDInFilters), arg0)
}

// ExistOriginalURLIDInFilters mocks base method.
func (m *MockCacheDAO) ExistOriginalURLIDInFilters(arg0 string) (bool, *business.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistOriginalURLIDInFilters", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*business.Error)
	return ret0, ret1
}

// ExistOriginalURLIDInFilters indicates an expected call of ExistOriginalURLIDInFilters.
func (mr *MockCacheDAOMockRecorder) ExistOriginalURLIDInFilters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistOriginalURLIDInFilters", reflect.TypeOf((*MockCacheDAO)(nil).ExistOriginalURLIDInFilters), arg0)
}

// GetOriginalURL mocks base method.
func (m *MockCacheDAO) GetOriginalURL(arg0 string) (string, *business.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURL", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*business.Error)
	return ret0, ret1
}

// GetOriginalURL indicates an expected call of GetOriginalURL.
func (mr *MockCacheDAOMockRecorder) GetOriginalURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURL", reflect.TypeOf((*MockCacheDAO)(nil).GetOriginalURL), arg0)
}

// SetOriginalURL mocks base method.
func (m *MockCacheDAO) SetOriginalURL(arg0, arg1 string) *business.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOriginalURL", arg0, arg1)
	ret0, _ := ret[0].(*business.Error)
	return ret0
}

// SetOriginalURL indicates an expected call of SetOriginalURL.
func (mr *MockCacheDAOMockRecorder) SetOriginalURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOriginalURL", reflect.TypeOf((*MockCacheDAO)(nil).SetOriginalURL), arg0, arg1)
}

// MockKeyDAO is a mock of KeyDAO interface.
type MockKeyDAO struct {
	ctrl     *gomock.Controller
	recorder *MockKeyDAOMockRecorder
}

// MockKeyDAOMockRecorder is the mock recorder for MockKeyDAO.
type MockKeyDAOMockRecorder struct {
	mock *MockKeyDAO
}

// NewMockKeyDAO creates a new mock instance.
func NewMockKeyDAO(ctrl *gomock.Controller) *MockKeyDAO {
	mock := &MockKeyDAO{ctrl: ctrl}
	mock.recorder = &MockKeyDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyDAO) EXPECT() *MockKeyDAOMockRecorder {
	return m.recorder
}

// BatchCreate mocks base method.
func (m *MockKeyDAO) BatchCreate(arg0 int) (int, *business.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchCreate", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*business.Error)
	return ret0, ret1
}

// BatchCreate indicates an expected call of BatchCreate.
func (mr *MockKeyDAOMockRecorder) BatchCreate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchCreate", reflect.TypeOf((*MockKeyDAO)(nil).BatchCreate), arg0)
}

// MockUrlDAO is a mock of UrlDAO interface.
type MockUrlDAO struct {
	ctrl     *gomock.Controller
	recorder *MockUrlDAOMockRecorder
}

// MockUrlDAOMockRecorder is the mock recorder for MockUrlDAO.
type MockUrlDAOMockRecorder struct {
	mock *MockUrlDAO
}

// NewMockUrlDAO creates a new mock instance.
func NewMockUrlDAO(ctrl *gomock.Controller) *MockUrlDAO {
	mock := &MockUrlDAO{ctrl: ctrl}
	mock.recorder = &MockUrlDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUrlDAO) EXPECT() *MockUrlDAOMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUrlDAO) Create(arg0 string) (*dao.URL, *business.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*dao.URL)
	ret1, _ := ret[1].(*business.Error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUrlDAOMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUrlDAO)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockUrlDAO) Delete(arg0 string) *business.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(*business.Error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUrlDAOMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUrlDAO)(nil).Delete), arg0)
}

// Expire mocks base method.
func (m *MockUrlDAO) Expire(arg0 int) ([]string, *business.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expire", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(*business.Error)
	return ret0, ret1
}

// Expire indicates an expected call of Expire.
func (mr *MockUrlDAOMockRecorder) Expire(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expire", reflect.TypeOf((*MockUrlDAO)(nil).Expire), arg0)
}

// Get mocks base method.
func (m *MockUrlDAO) Get(arg0 string) (*dao.URL, *business.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*dao.URL)
	ret1, _ := ret[1].(*business.Error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUrlDAOMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUrlDAO)(nil).Get), arg0)
}
