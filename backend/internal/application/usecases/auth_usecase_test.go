package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/atdevten/peace/internal/application/commands"
	"github.com/atdevten/peace/internal/application/services/google"
	appjwt "github.com/atdevten/peace/internal/application/services/jwt"
	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/value_objects"
	"github.com/atdevten/peace/testutils/helpers"
	repositories "github.com/atdevten/peace/testutils/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Mock services
type MockJWTService struct {
	ctrl                     *gomock.Controller
	ValidateRefreshTokenFunc func(token string) (*appjwt.Claims, error)
}

func (m *MockJWTService) GenerateAccessToken(userID value_objects.UserID, email value_objects.Email) (string, error) {
	return "mock-access-token", nil
}

func (m *MockJWTService) GenerateRefreshToken(userID value_objects.UserID, email value_objects.Email) (string, error) {
	return "mock-refresh-token", nil
}

func (m *MockJWTService) ValidateAccessToken(token string) (*appjwt.Claims, error) {
	return &appjwt.Claims{
		UserID: "550e8400-e29b-41d4-a716-446655440000",
		Email:  "test@example.com",
	}, nil
}

func (m *MockJWTService) ValidateRefreshToken(token string) (*appjwt.Claims, error) {
	if m.ValidateRefreshTokenFunc != nil {
		return m.ValidateRefreshTokenFunc(token)
	}
	return &appjwt.Claims{
		UserID: "550e8400-e29b-41d4-a716-446655440000",
		Email:  "test@example.com",
	}, nil
}

type MockGoogleService struct {
	ctrl                     *gomock.Controller
	ExchangeCodeForTokenFunc func(ctx context.Context, code string) (*google.GoogleUserInfo, error)
}

func (m *MockGoogleService) GetAuthURL() string {
	return "https://accounts.google.com/o/oauth2/v2/auth?mock=true"
}

func (m *MockGoogleService) ExchangeCodeForToken(ctx context.Context, code string) (*google.GoogleUserInfo, error) {
	if m.ExchangeCodeForTokenFunc != nil {
		return m.ExchangeCodeForTokenFunc(ctx, code)
	}
	return &google.GoogleUserInfo{
		ID:        "google123",
		Email:     "google@example.com",
		FirstName: "Google",
		LastName:  "User",
		Picture:   "https://example.com/avatar.jpg",
	}, nil
}

func TestAuthUseCaseImpl_Register(t *testing.T) {
	tests := []struct {
		name        string
		command     commands.RegisterCommand
		mockUser    *entities.User
		emailExists bool
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful registration",
			command: commands.RegisterCommand{
				Email:     "test@example.com",
				Username:  "testuser",
				FirstName: helpers.StringPtr("John"),
				LastName:  helpers.StringPtr("Doe"),
				Password:  "Password123",
			},
			mockUser:    helpers.CreateTestUser(),
			emailExists: false,
			wantErr:     false,
		},
		{
			name: "email already exists",
			command: commands.RegisterCommand{
				Email:     "existing@example.com",
				Username:  "testuser",
				FirstName: helpers.StringPtr("John"),
				LastName:  helpers.StringPtr("Doe"),
				Password:  "Password123",
			},
			mockUser:    helpers.CreateTestUser(),
			emailExists: true,
			wantErr:     true,
			expectedErr: "email already registered",
		},
		{
			name: "invalid email format",
			command: commands.RegisterCommand{
				Email:     "invalid-email",
				Username:  "testuser",
				FirstName: helpers.StringPtr("John"),
				LastName:  helpers.StringPtr("Doe"),
				Password:  "Password123",
			},
			wantErr:     true,
			expectedErr: "invalid email format",
		},
		{
			name: "invalid password",
			command: commands.RegisterCommand{
				Email:     "test@example.com",
				Username:  "testuser",
				FirstName: helpers.StringPtr("John"),
				LastName:  helpers.StringPtr("Doe"),
				Password:  "weak",
			},
			wantErr:     true,
			expectedErr: "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)
			if tt.command.Email != "invalid-email" && tt.command.Password != "weak" {
				mockRepo.EXPECT().EmailExists(gomock.Any(), gomock.Any()).Return(tt.emailExists, tt.mockError)
				if !tt.emailExists {
					mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			// Setup mock services
			mockJWT := &MockJWTService{ctrl: ctrl}
			mockGoogle := &MockGoogleService{ctrl: ctrl}

			useCase := NewAuthUseCase(mockRepo, mockJWT, mockGoogle)
			user, err := useCase.Register(context.Background(), tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.command.Email, user.Email().String())
			}
		})
	}
}

func TestAuthUseCaseImpl_Login(t *testing.T) {
	tests := []struct {
		name        string
		command     commands.LoginCommand
		mockUser    *entities.User
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful login",
			command: commands.LoginCommand{
				Email:    "test@example.com",
				Password: "Password123",
			},
			mockUser: helpers.CreateTestUser(),
			wantErr:  false,
		},
		{
			name: "invalid email format",
			command: commands.LoginCommand{
				Email:    "invalid-email",
				Password: "Password123",
			},
			wantErr:     true,
			expectedErr: "invalid email format",
		},
		{
			name: "user not found",
			command: commands.LoginCommand{
				Email:    "notfound@example.com",
				Password: "Password123",
			},
			mockError:   errors.New("user not found"),
			wantErr:     true,
			expectedErr: "invalid email or password",
		},
		{
			name: "wrong password",
			command: commands.LoginCommand{
				Email:    "test@example.com",
				Password: "WrongPassword",
			},
			mockUser:    helpers.CreateTestUser(),
			wantErr:     true,
			expectedErr: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)
			if tt.command.Email != "invalid-email" {
				mockRepo.EXPECT().GetByFilter(gomock.Any(), gomock.Any()).Return(tt.mockUser, tt.mockError)
			}

			// Setup mock services
			mockJWT := &MockJWTService{ctrl: ctrl}
			mockGoogle := &MockGoogleService{ctrl: ctrl}

			useCase := NewAuthUseCase(mockRepo, mockJWT, mockGoogle)
			user, access, refresh, err := useCase.Login(context.Background(), tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
				assert.Empty(t, access)
				assert.Empty(t, refresh)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "mock-access-token", access)
				assert.Equal(t, "mock-refresh-token", refresh)
			}
		})
	}
}

func TestAuthUseCaseImpl_Refresh(t *testing.T) {
	tests := []struct {
		name         string
		accessToken  string
		refreshToken string
		mockError    error
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "successful refresh",
			accessToken:  "valid-access-token",
			refreshToken: "valid-refresh-token",
			wantErr:      false,
		},
		{
			name:         "invalid refresh token",
			accessToken:  "valid-access-token",
			refreshToken: "invalid-refresh-token",
			mockError:    errors.New("invalid token"),
			wantErr:      true,
			expectedErr:  "invalid refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)

			// Setup mock services
			mockJWT := &MockJWTService{ctrl: ctrl}
			if tt.mockError != nil {
				mockJWT.ValidateRefreshTokenFunc = func(token string) (*appjwt.Claims, error) {
					return nil, tt.mockError
				}
			}
			mockGoogle := &MockGoogleService{ctrl: ctrl}

			useCase := NewAuthUseCase(mockRepo, mockJWT, mockGoogle)
			newAccess, newRefresh, err := useCase.Refresh(context.Background(), tt.accessToken, tt.refreshToken)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, newAccess)
				assert.Empty(t, newRefresh)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "mock-access-token", newAccess)
				assert.Equal(t, "mock-refresh-token", newRefresh)
			}
		})
	}
}

func TestAuthUseCaseImpl_LoginWithGoogle(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		mockUser    *entities.User
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "successful Google login - existing user",
			code:     "valid-google-code",
			mockUser: helpers.CreateTestGoogleUser(),
			wantErr:  false,
		},
		{
			name:      "successful Google login - new user",
			code:      "valid-google-code",
			mockError: errors.New("user not found"),
			wantErr:   false,
		},
		{
			name:        "invalid Google code",
			code:        "invalid-google-code",
			mockError:   errors.New("invalid code"),
			wantErr:     true,
			expectedErr: "invalid code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)
			if tt.code == "valid-google-code" {
				mockRepo.EXPECT().GetByFilter(gomock.Any(), gomock.Any()).Return(tt.mockUser, tt.mockError)
				if tt.mockError != nil {
					// New user case - expect Create call
					mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			// Setup mock services
			mockJWT := &MockJWTService{ctrl: ctrl}
			mockGoogle := &MockGoogleService{ctrl: ctrl}
			if tt.mockError != nil && tt.code == "invalid-google-code" {
				mockGoogle.ExchangeCodeForTokenFunc = func(ctx context.Context, code string) (*google.GoogleUserInfo, error) {
					return nil, tt.mockError
				}
			}

			useCase := NewAuthUseCase(mockRepo, mockJWT, mockGoogle)
			user, access, refresh, err := useCase.LoginWithGoogle(context.Background(), tt.code)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
				assert.Empty(t, access)
				assert.Empty(t, refresh)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "mock-access-token", access)
				assert.Equal(t, "mock-refresh-token", refresh)
			}
		})
	}
}
