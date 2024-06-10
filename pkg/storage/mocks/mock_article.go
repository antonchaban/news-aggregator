// Code generated by MockGen. DO NOT EDIT.
// Source: news-aggregator/pkg/storage (interfaces: ArticleStorage)

// Package mocks is a generated GoMock package.
package mocks

import (
	model "news-aggregator/pkg/model"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockArticleStorage is a mock of ArticleStorage interface.
type MockArticleStorage struct {
	ctrl     *gomock.Controller
	recorder *MockArticleStorageMockRecorder
}

// MockArticleStorageMockRecorder is the mock recorder for MockArticleStorage.
type MockArticleStorageMockRecorder struct {
	mock *MockArticleStorage
}

// NewMockArticleStorage creates a new mock instance.
func NewMockArticleStorage(ctrl *gomock.Controller) *MockArticleStorage {
	mock := &MockArticleStorage{ctrl: ctrl}
	mock.recorder = &MockArticleStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleStorage) EXPECT() *MockArticleStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockArticleStorage) Create(arg0 model.Article) (model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockArticleStorageMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockArticleStorage)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockArticleStorage) Delete(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockArticleStorageMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockArticleStorage)(nil).Delete), arg0)
}

// GetAll mocks base method.
func (m *MockArticleStorage) GetAll() ([]model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockArticleStorageMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockArticleStorage)(nil).GetAll))
}

// GetByDateInRange mocks base method.
func (m *MockArticleStorage) GetByDateInRange(arg0, arg1 time.Time) ([]model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByDateInRange", arg0, arg1)
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByDateInRange indicates an expected call of GetByDateInRange.
func (mr *MockArticleStorageMockRecorder) GetByDateInRange(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByDateInRange", reflect.TypeOf((*MockArticleStorage)(nil).GetByDateInRange), arg0, arg1)
}

// GetByKeyword mocks base method.
func (m *MockArticleStorage) GetByKeyword(arg0 string) ([]model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKeyword", arg0)
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKeyword indicates an expected call of GetByKeyword.
func (mr *MockArticleStorageMockRecorder) GetByKeyword(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKeyword", reflect.TypeOf((*MockArticleStorage)(nil).GetByKeyword), arg0)
}

// GetBySource mocks base method.
func (m *MockArticleStorage) GetBySource(arg0 string) ([]model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySource", arg0)
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBySource indicates an expected call of GetBySource.
func (mr *MockArticleStorageMockRecorder) GetBySource(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySource", reflect.TypeOf((*MockArticleStorage)(nil).GetBySource), arg0)
}

// SaveAll mocks base method.
func (m *MockArticleStorage) SaveAll(arg0 []model.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAll indicates an expected call of SaveAll.
func (mr *MockArticleStorageMockRecorder) SaveAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAll", reflect.TypeOf((*MockArticleStorage)(nil).SaveAll), arg0)
}