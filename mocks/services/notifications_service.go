// Code generated by MockGen. DO NOT EDIT.
// Source: services.go
//
// Generated by this command:
//
//	mockgen -source=services.go -destination=../../mocks/services/notifications_service.go -package=mockservices -exclude_interfaces=ToysService,FileStorageService,TicketsService,SsoService
//

// Package mockservices is a generated GoMock package.
package mockservices

import (
	context "context"
	reflect "reflect"

	entities "github.com/DKhorkov/hmtm-bff/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockNotificationsService is a mock of NotificationsService interface.
type MockNotificationsService struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationsServiceMockRecorder
	isgomock struct{}
}

// MockNotificationsServiceMockRecorder is the mock recorder for MockNotificationsService.
type MockNotificationsServiceMockRecorder struct {
	mock *MockNotificationsService
}

// NewMockNotificationsService creates a new mock instance.
func NewMockNotificationsService(ctrl *gomock.Controller) *MockNotificationsService {
	mock := &MockNotificationsService{ctrl: ctrl}
	mock.recorder = &MockNotificationsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationsService) EXPECT() *MockNotificationsServiceMockRecorder {
	return m.recorder
}

// CountUserEmailCommunications mocks base method.
func (m *MockNotificationsService) CountUserEmailCommunications(ctx context.Context, userID uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUserEmailCommunications", ctx, userID)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUserEmailCommunications indicates an expected call of CountUserEmailCommunications.
func (mr *MockNotificationsServiceMockRecorder) CountUserEmailCommunications(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUserEmailCommunications", reflect.TypeOf((*MockNotificationsService)(nil).CountUserEmailCommunications), ctx, userID)
}

// GetUserEmailCommunications mocks base method.
func (m *MockNotificationsService) GetUserEmailCommunications(ctx context.Context, userID uint64, pagination *entities.Pagination) ([]entities.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserEmailCommunications", ctx, userID, pagination)
	ret0, _ := ret[0].([]entities.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserEmailCommunications indicates an expected call of GetUserEmailCommunications.
func (mr *MockNotificationsServiceMockRecorder) GetUserEmailCommunications(ctx, userID, pagination any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserEmailCommunications", reflect.TypeOf((*MockNotificationsService)(nil).GetUserEmailCommunications), ctx, userID, pagination)
}
