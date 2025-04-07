// Code generated by MockGen. DO NOT EDIT.
// Source: repositories.go
//
// Generated by this command:
//
//	mockgen -source=repositories.go -destination=../../mocks/repositories/sso_repository.go -package=mockrepositories -exclude_interfaces=ToysRepository,FileStorageRepository,TicketsRepository,NotificationsRepository
//

// Package mockrepositories is a generated GoMock package.
package mockrepositories

import (
	context "context"
	reflect "reflect"

	entities "github.com/DKhorkov/hmtm-bff/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockSsoRepository is a mock of SsoRepository interface.
type MockSsoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSsoRepositoryMockRecorder
	isgomock struct{}
}

// MockSsoRepositoryMockRecorder is the mock recorder for MockSsoRepository.
type MockSsoRepositoryMockRecorder struct {
	mock *MockSsoRepository
}

// NewMockSsoRepository creates a new mock instance.
func NewMockSsoRepository(ctrl *gomock.Controller) *MockSsoRepository {
	mock := &MockSsoRepository{ctrl: ctrl}
	mock.recorder = &MockSsoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSsoRepository) EXPECT() *MockSsoRepositoryMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockSsoRepository) ChangePassword(ctx context.Context, accessToken, oldPassword, newPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", ctx, accessToken, oldPassword, newPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockSsoRepositoryMockRecorder) ChangePassword(ctx, accessToken, oldPassword, newPassword any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockSsoRepository)(nil).ChangePassword), ctx, accessToken, oldPassword, newPassword)
}

// ForgetPassword mocks base method.
func (m *MockSsoRepository) ForgetPassword(ctx context.Context, forgetPasswordToken, newPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForgetPassword", ctx, forgetPasswordToken, newPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForgetPassword indicates an expected call of ForgetPassword.
func (mr *MockSsoRepositoryMockRecorder) ForgetPassword(ctx, forgetPasswordToken, newPassword any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgetPassword", reflect.TypeOf((*MockSsoRepository)(nil).ForgetPassword), ctx, forgetPasswordToken, newPassword)
}

// GetAllUsers mocks base method.
func (m *MockSsoRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", ctx)
	ret0, _ := ret[0].([]entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockSsoRepositoryMockRecorder) GetAllUsers(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockSsoRepository)(nil).GetAllUsers), ctx)
}

// GetMe mocks base method.
func (m *MockSsoRepository) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMe", ctx, accessToken)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMe indicates an expected call of GetMe.
func (mr *MockSsoRepositoryMockRecorder) GetMe(ctx, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMe", reflect.TypeOf((*MockSsoRepository)(nil).GetMe), ctx, accessToken)
}

// GetUserByEmail mocks base method.
func (m *MockSsoRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockSsoRepositoryMockRecorder) GetUserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockSsoRepository)(nil).GetUserByEmail), ctx, email)
}

// GetUserByID mocks base method.
func (m *MockSsoRepository) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockSsoRepositoryMockRecorder) GetUserByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockSsoRepository)(nil).GetUserByID), ctx, id)
}

// LoginUser mocks base method.
func (m *MockSsoRepository) LoginUser(ctx context.Context, userData entities.LoginUserDTO) (*entities.TokensDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, userData)
	ret0, _ := ret[0].(*entities.TokensDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockSsoRepositoryMockRecorder) LoginUser(ctx, userData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockSsoRepository)(nil).LoginUser), ctx, userData)
}

// LogoutUser mocks base method.
func (m *MockSsoRepository) LogoutUser(ctx context.Context, accessToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoutUser", ctx, accessToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogoutUser indicates an expected call of LogoutUser.
func (mr *MockSsoRepositoryMockRecorder) LogoutUser(ctx, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoutUser", reflect.TypeOf((*MockSsoRepository)(nil).LogoutUser), ctx, accessToken)
}

// RefreshTokens mocks base method.
func (m *MockSsoRepository) RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshTokens", ctx, refreshToken)
	ret0, _ := ret[0].(*entities.TokensDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshTokens indicates an expected call of RefreshTokens.
func (mr *MockSsoRepositoryMockRecorder) RefreshTokens(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTokens", reflect.TypeOf((*MockSsoRepository)(nil).RefreshTokens), ctx, refreshToken)
}

// RegisterUser mocks base method.
func (m *MockSsoRepository) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, userData)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockSsoRepositoryMockRecorder) RegisterUser(ctx, userData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockSsoRepository)(nil).RegisterUser), ctx, userData)
}

// SendForgetPasswordMessage mocks base method.
func (m *MockSsoRepository) SendForgetPasswordMessage(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendForgetPasswordMessage", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendForgetPasswordMessage indicates an expected call of SendForgetPasswordMessage.
func (mr *MockSsoRepositoryMockRecorder) SendForgetPasswordMessage(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendForgetPasswordMessage", reflect.TypeOf((*MockSsoRepository)(nil).SendForgetPasswordMessage), ctx, email)
}

// SendVerifyEmailMessage mocks base method.
func (m *MockSsoRepository) SendVerifyEmailMessage(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendVerifyEmailMessage", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendVerifyEmailMessage indicates an expected call of SendVerifyEmailMessage.
func (mr *MockSsoRepositoryMockRecorder) SendVerifyEmailMessage(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendVerifyEmailMessage", reflect.TypeOf((*MockSsoRepository)(nil).SendVerifyEmailMessage), ctx, email)
}

// UpdateUserProfile mocks base method.
func (m *MockSsoRepository) UpdateUserProfile(ctx context.Context, userProfileData entities.UpdateUserProfileDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", ctx, userProfileData)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockSsoRepositoryMockRecorder) UpdateUserProfile(ctx, userProfileData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockSsoRepository)(nil).UpdateUserProfile), ctx, userProfileData)
}

// VerifyUserEmail mocks base method.
func (m *MockSsoRepository) VerifyUserEmail(ctx context.Context, verifyEmailToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyUserEmail", ctx, verifyEmailToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyUserEmail indicates an expected call of VerifyUserEmail.
func (mr *MockSsoRepositoryMockRecorder) VerifyUserEmail(ctx, verifyEmailToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyUserEmail", reflect.TypeOf((*MockSsoRepository)(nil).VerifyUserEmail), ctx, verifyEmailToken)
}
