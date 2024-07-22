// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/antonchaban/news-aggregator/pkg/service (interfaces: SourceService)
//
// Generated by this command:
//
//	mockgen -destination=../service/mocks/mock_source_service.go -package=mocks github.com/antonchaban/news-aggregator/pkg/service SourceService
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/antonchaban/news-aggregator/pkg/model"
	gomock "go.uber.org/mock/gomock"
)

// MockSourceService is a mock of SourceService interface.
type MockSourceService struct {
	ctrl     *gomock.Controller
	recorder *MockSourceServiceMockRecorder
}

// MockSourceServiceMockRecorder is the mock recorder for MockSourceService.
type MockSourceServiceMockRecorder struct {
	mock *MockSourceService
}

// NewMockSourceService creates a new mock instance.
func NewMockSourceService(ctrl *gomock.Controller) *MockSourceService {
	mock := &MockSourceService{ctrl: ctrl}
	mock.recorder = &MockSourceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSourceService) EXPECT() *MockSourceServiceMockRecorder {
	return m.recorder
}

// AddSource mocks base method.
func (m *MockSourceService) AddSource(arg0 model.Source) (model.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSource", arg0)
	ret0, _ := ret[0].(model.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSource indicates an expected call of AddSource.
func (mr *MockSourceServiceMockRecorder) AddSource(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSource", reflect.TypeOf((*MockSourceService)(nil).AddSource), arg0)
}

// DeleteSource mocks base method.
func (m *MockSourceService) DeleteSource(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSource", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSource indicates an expected call of DeleteSource.
func (mr *MockSourceServiceMockRecorder) DeleteSource(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSource", reflect.TypeOf((*MockSourceService)(nil).DeleteSource), arg0)
}

// FetchFromAllSources mocks base method.
func (m *MockSourceService) FetchFromAllSources() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchFromAllSources")
	ret0, _ := ret[0].(error)
	return ret0
}

// FetchFromAllSources indicates an expected call of FetchFromAllSources.
func (mr *MockSourceServiceMockRecorder) FetchFromAllSources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchFromAllSources", reflect.TypeOf((*MockSourceService)(nil).FetchFromAllSources))
}

// FetchSourceByID mocks base method.
func (m *MockSourceService) FetchSourceByID(arg0 int) ([]model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchSourceByID", arg0)
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchSourceByID indicates an expected call of FetchSourceByID.
func (mr *MockSourceServiceMockRecorder) FetchSourceByID(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchSourceByID", reflect.TypeOf((*MockSourceService)(nil).FetchSourceByID), arg0)
}

// GetAll mocks base method.
func (m *MockSourceService) GetAll() ([]model.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]model.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockSourceServiceMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockSourceService)(nil).GetAll))
}

// LoadDataFromFiles mocks base method.
func (m *MockSourceService) LoadDataFromFiles() ([]model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadDataFromFiles")
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadDataFromFiles indicates an expected call of LoadDataFromFiles.
func (mr *MockSourceServiceMockRecorder) LoadDataFromFiles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadDataFromFiles", reflect.TypeOf((*MockSourceService)(nil).LoadDataFromFiles))
}

// UpdateSource mocks base method.
func (m *MockSourceService) UpdateSource(arg0 int, arg1 model.Source) (model.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSource", arg0, arg1)
	ret0, _ := ret[0].(model.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSource indicates an expected call of UpdateSource.
func (mr *MockSourceServiceMockRecorder) UpdateSource(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSource", reflect.TypeOf((*MockSourceService)(nil).UpdateSource), arg0, arg1)
}
