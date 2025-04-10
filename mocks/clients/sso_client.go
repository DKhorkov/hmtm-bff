// Code generated by MockGen. DO NOT EDIT.
// Source: clients.go
//
// Generated by this command:
//
//	mockgen -source=clients.go -destination=../../mocks/clients/sso_client.go -package=mockclients -exclude_interfaces=ToysClient,TicketsClient,NotificationsClient,S3Client
//

// Package mockclients is a generated GoMock package.
package mockclients

import (
	context "context"
	reflect "reflect"

	sso "github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockSsoClient is a mock of SsoClient interface.
type MockSsoClient struct {
	ctrl     *gomock.Controller
	recorder *MockSsoClientMockRecorder
	isgomock struct{}
}

// MockSsoClientMockRecorder is the mock recorder for MockSsoClient.
type MockSsoClientMockRecorder struct {
	mock *MockSsoClient
}

// NewMockSsoClient creates a new mock instance.
func NewMockSsoClient(ctrl *gomock.Controller) *MockSsoClient {
	mock := &MockSsoClient{ctrl: ctrl}
	mock.recorder = &MockSsoClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSsoClient) EXPECT() *MockSsoClientMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockSsoClient) ChangePassword(ctx context.Context, in *sso.ChangePasswordIn, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ChangePassword", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockSsoClientMockRecorder) ChangePassword(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockSsoClient)(nil).ChangePassword), varargs...)
}

// ForgetPassword mocks base method.
func (m *MockSsoClient) ForgetPassword(ctx context.Context, in *sso.ForgetPasswordIn, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ForgetPassword", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForgetPassword indicates an expected call of ForgetPassword.
func (mr *MockSsoClientMockRecorder) ForgetPassword(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgetPassword", reflect.TypeOf((*MockSsoClient)(nil).ForgetPassword), varargs...)
}

// GetMe mocks base method.
func (m *MockSsoClient) GetMe(ctx context.Context, in *sso.GetMeIn, opts ...grpc.CallOption) (*sso.GetUserOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMe", varargs...)
	ret0, _ := ret[0].(*sso.GetUserOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMe indicates an expected call of GetMe.
func (mr *MockSsoClientMockRecorder) GetMe(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMe", reflect.TypeOf((*MockSsoClient)(nil).GetMe), varargs...)
}

// GetUser mocks base method.
func (m *MockSsoClient) GetUser(ctx context.Context, in *sso.GetUserIn, opts ...grpc.CallOption) (*sso.GetUserOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUser", varargs...)
	ret0, _ := ret[0].(*sso.GetUserOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockSsoClientMockRecorder) GetUser(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockSsoClient)(nil).GetUser), varargs...)
}

// GetUserByEmail mocks base method.
func (m *MockSsoClient) GetUserByEmail(ctx context.Context, in *sso.GetUserByEmailIn, opts ...grpc.CallOption) (*sso.GetUserOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserByEmail", varargs...)
	ret0, _ := ret[0].(*sso.GetUserOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockSsoClientMockRecorder) GetUserByEmail(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockSsoClient)(nil).GetUserByEmail), varargs...)
}

// GetUsers mocks base method.
func (m *MockSsoClient) GetUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*sso.GetUsersOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUsers", varargs...)
	ret0, _ := ret[0].(*sso.GetUsersOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockSsoClientMockRecorder) GetUsers(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockSsoClient)(nil).GetUsers), varargs...)
}

// Login mocks base method.
func (m *MockSsoClient) Login(ctx context.Context, in *sso.LoginIn, opts ...grpc.CallOption) (*sso.LoginOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Login", varargs...)
	ret0, _ := ret[0].(*sso.LoginOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockSsoClientMockRecorder) Login(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockSsoClient)(nil).Login), varargs...)
}

// Logout mocks base method.
func (m *MockSsoClient) Logout(ctx context.Context, in *sso.LogoutIn, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Logout", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Logout indicates an expected call of Logout.
func (mr *MockSsoClientMockRecorder) Logout(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockSsoClient)(nil).Logout), varargs...)
}

// RefreshTokens mocks base method.
func (m *MockSsoClient) RefreshTokens(ctx context.Context, in *sso.RefreshTokensIn, opts ...grpc.CallOption) (*sso.LoginOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RefreshTokens", varargs...)
	ret0, _ := ret[0].(*sso.LoginOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshTokens indicates an expected call of RefreshTokens.
func (mr *MockSsoClientMockRecorder) RefreshTokens(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTokens", reflect.TypeOf((*MockSsoClient)(nil).RefreshTokens), varargs...)
}

// Register mocks base method.
func (m *MockSsoClient) Register(ctx context.Context, in *sso.RegisterIn, opts ...grpc.CallOption) (*sso.RegisterOut, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Register", varargs...)
	ret0, _ := ret[0].(*sso.RegisterOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockSsoClientMockRecorder) Register(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockSsoClient)(nil).Register), varargs...)
}

// SendForgetPasswordMessage mocks base method.
func (m *MockSsoClient) SendForgetPasswordMessage(ctx context.Context, in *sso.SendForgetPasswordMessageIn, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendForgetPasswordMessage", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendForgetPasswordMessage indicates an expected call of SendForgetPasswordMessage.
func (mr *MockSsoClientMockRecorder) SendForgetPasswordMessage(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendForgetPasswordMessage", reflect.TypeOf((*MockSsoClient)(nil).SendForgetPasswordMessage), varargs...)
}

// SendVerifyEmailMessage mocks base method.
func (m *MockSsoClient) SendVerifyEmailMessage(ctx context.Context, in *sso.SendVerifyEmailMessageIn, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendVerifyEmailMessage", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendVerifyEmailMessage indicates an expected call of SendVerifyEmailMessage.
func (mr *MockSsoClientMockRecorder) SendVerifyEmailMessage(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendVerifyEmailMessage", reflect.TypeOf((*MockSsoClient)(nil).SendVerifyEmailMessage), varargs...)
}

// UpdateUserProfile mocks base method.
func (m *MockSsoClient) UpdateUserProfile(ctx context.Context, in *sso.UpdateUserProfileIn, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateUserProfile", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockSsoClientMockRecorder) UpdateUserProfile(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockSsoClient)(nil).UpdateUserProfile), varargs...)
}

// VerifyEmail mocks base method.
func (m *MockSsoClient) VerifyEmail(ctx context.Context, in *sso.VerifyEmailIn, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "VerifyEmail", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyEmail indicates an expected call of VerifyEmail.
func (mr *MockSsoClientMockRecorder) VerifyEmail(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyEmail", reflect.TypeOf((*MockSsoClient)(nil).VerifyEmail), varargs...)
}
