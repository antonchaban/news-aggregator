// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/antonchaban/news-aggregator/pkg/storage (interfaces: ArticleStorage)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/antonchaban/news-aggregator/pkg/model"
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

// DeleteBySourceID mocks base method.
func (m *MockArticleStorage) DeleteBySourceID(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBySourceID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBySourceID indicates an expected call of DeleteBySourceID.
func (mr *MockArticleStorageMockRecorder) DeleteBySourceID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBySourceID", reflect.TypeOf((*MockArticleStorage)(nil).DeleteBySourceID), arg0)
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

// GetByFilter mocks base method.
func (m *MockArticleStorage) GetByFilter(arg0 string, arg1 []interface{}) ([]model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByFilter", arg0, arg1)
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByFilter indicates an expected call of GetByFilter.
func (mr *MockArticleStorageMockRecorder) GetByFilter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByFilter", reflect.TypeOf((*MockArticleStorage)(nil).GetByFilter), arg0, arg1)
}

// Save mocks base method.
func (m *MockArticleStorage) Save(arg0 model.Article) (model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockArticleStorageMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockArticleStorage)(nil).Save), arg0)
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
