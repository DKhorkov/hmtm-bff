// Code generated by MockGen. DO NOT EDIT.
// Source: services.go
//
// Generated by this command:
//
//	mockgen -source=services.go -destination=../../mocks/services/file_storage_service.go -package=mockservices -exclude_interfaces=ToysService,SsoService,TicketsService,NotificationsService
//

// Package mockservices is a generated GoMock package.
package mockservices

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFileStorageService is a mock of FileStorageService interface.
type MockFileStorageService struct {
	ctrl     *gomock.Controller
	recorder *MockFileStorageServiceMockRecorder
	isgomock struct{}
}

// MockFileStorageServiceMockRecorder is the mock recorder for MockFileStorageService.
type MockFileStorageServiceMockRecorder struct {
	mock *MockFileStorageService
}

// NewMockFileStorageService creates a new mock instance.
func NewMockFileStorageService(ctrl *gomock.Controller) *MockFileStorageService {
	mock := &MockFileStorageService{ctrl: ctrl}
	mock.recorder = &MockFileStorageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileStorageService) EXPECT() *MockFileStorageServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockFileStorageService) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockFileStorageServiceMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFileStorageService)(nil).Delete), ctx, key)
}

// DeleteMany mocks base method.
func (m *MockFileStorageService) DeleteMany(ctx context.Context, keys []string) []error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMany", ctx, keys)
	ret0, _ := ret[0].([]error)
	return ret0
}

// DeleteMany indicates an expected call of DeleteMany.
func (mr *MockFileStorageServiceMockRecorder) DeleteMany(ctx, keys any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMany", reflect.TypeOf((*MockFileStorageService)(nil).DeleteMany), ctx, keys)
}

// Upload mocks base method.
func (m *MockFileStorageService) Upload(ctx context.Context, key string, file []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, key, file)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockFileStorageServiceMockRecorder) Upload(ctx, key, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockFileStorageService)(nil).Upload), ctx, key, file)
}
