// Code generated by MockGen. DO NOT EDIT.
// Source: clients.go
//
// Generated by this command:
//
//	mockgen -source=clients.go -destination=../../mocks/clients/notifications_client.go -package=mockclients -exclude_interfaces=ToysClient,TicketsClient,SsoClient,S3Client
//

// Package mockclients is a generated GoMock package.
package mockclients

import (
	context "context"
	reflect "reflect"

	notifications "github.com/DKhorkov/hmtm-notifications/api/protobuf/generated/go/notifications"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockNotificationsClient is a mock of NotificationsClient interface.
type MockNotificationsClient struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationsClientMockRecorder
	isgomock struct{}
}

// MockNotificationsClientMockRecorder is the mock recorder for MockNotificationsClient.
type MockNotificationsClientMockRecorder struct {
	mock *MockNotificationsClient
}

// NewMockNotificationsClient creates a new mock instance.
func NewMockNotificationsClient(ctrl *gomock.Controller) *MockNotificationsClient {
	mock := &MockNotificationsClient{ctrl: ctrl}
	mock.recorder = &MockNotificationsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationsClient) EXPECT() *MockNotificationsClientMockRecorder {
	return m.recorder
}

// CountUserEmailCommunications mocks base method.
func (m *MockNotificationsClient) CountUserEmailCommunications(ctx context.Context, in *notifications.CountUserEmailCommunicationsIn, opts ...grpc.CallOption) (*notifications.CountOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CountUserEmailCommunications", varargs...)
	ret0, _ := ret[0].(*notifications.CountOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUserEmailCommunications indicates an expected call of CountUserEmailCommunications.
func (mr *MockNotificationsClientMockRecorder) CountUserEmailCommunications(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUserEmailCommunications", reflect.TypeOf((*MockNotificationsClient)(nil).CountUserEmailCommunications), varargs...)
}

// GetUserEmailCommunications mocks base method.
func (m *MockNotificationsClient) GetUserEmailCommunications(ctx context.Context, in *notifications.GetUserEmailCommunicationsIn, opts ...grpc.CallOption) (*notifications.GetUserEmailCommunicationsOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserEmailCommunications", varargs...)
	ret0, _ := ret[0].(*notifications.GetUserEmailCommunicationsOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserEmailCommunications indicates an expected call of GetUserEmailCommunications.
func (mr *MockNotificationsClientMockRecorder) GetUserEmailCommunications(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserEmailCommunications", reflect.TypeOf((*MockNotificationsClient)(nil).GetUserEmailCommunications), varargs...)
}
