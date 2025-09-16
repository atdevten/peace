package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/testutils/helpers"
	repositories "github.com/atdevten/peace/testutils/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// CreateTestUserOnlineStatus creates a test user online status
func CreateTestUserOnlineStatus() *entities.UserOnlineStatus {
	userID := helpers.CreateTestUserID()
	email := helpers.CreateTestEmail("test@example.com")
	status := entities.NewUserOnlineStatus(userID, email)
	return status
}

func TestUserOnlineStatusUseCaseImpl_SetUserOnline(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		userEmail   string
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:      "successful set user online",
			userID:    "550e8400-e29b-41d4-a716-446655440000",
			userEmail: "test@example.com",
			wantErr:   false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			userEmail:   "test@example.com",
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "invalid email",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			userEmail:   "invalid-email",
			wantErr:     true,
			expectedErr: "invalid email format",
		},
		{
			name:        "repository error",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			userEmail:   "test@example.com",
			mockError:   errors.New("database error"),
			wantErr:     true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserOnlineStatusRepository(ctrl)
			if tt.userID != "invalid-id" && tt.userEmail != "invalid-email" {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(tt.mockError)
			}

			useCase := NewUserOnlineStatusUseCase(mockRepo)
			err := useCase.SetUserOnline(context.Background(), tt.userID, tt.userEmail)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserOnlineStatusUseCaseImpl_SetUserOffline(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful set user offline",
			userID:  "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			mockError:   errors.New("invalid UUID length"),
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "repository error",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockError:   errors.New("database error"),
			wantErr:     true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserOnlineStatusRepository(ctrl)
			mockRepo.EXPECT().GetByUserID(gomock.Any(), gomock.Any()).Return(CreateTestUserOnlineStatus(), nil)
			mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(tt.mockError)

			useCase := NewUserOnlineStatusUseCase(mockRepo)
			err := useCase.SetUserOffline(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserOnlineStatusUseCaseImpl_UpdateUserLastSeen(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful update last seen",
			userID:  "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			mockError:   errors.New("invalid UUID length"),
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "repository error",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockError:   errors.New("database error"),
			wantErr:     true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserOnlineStatusRepository(ctrl)
			mockRepo.EXPECT().UpdateLastSeen(gomock.Any(), gomock.Any()).Return(tt.mockError)

			useCase := NewUserOnlineStatusUseCase(mockRepo)
			err := useCase.UpdateUserLastSeen(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserOnlineStatusUseCaseImpl_GetUserOnlineStatus(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockStatus  *entities.UserOnlineStatus
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:       "successful get user online status",
			userID:     "550e8400-e29b-41d4-a716-446655440000",
			mockStatus: CreateTestUserOnlineStatus(),
			wantErr:    false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			mockError:   errors.New("invalid UUID length"),
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "status not found",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockError:   errors.New("status not found"),
			wantErr:     true,
			expectedErr: "status not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserOnlineStatusRepository(ctrl)
			mockRepo.EXPECT().GetByUserID(gomock.Any(), gomock.Any()).Return(tt.mockStatus, tt.mockError)

			useCase := NewUserOnlineStatusUseCase(mockRepo)
			status, err := useCase.GetUserOnlineStatus(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, status)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, status)
			}
		})
	}
}

func TestUserOnlineStatusUseCaseImpl_GetOnlineUsers(t *testing.T) {
	tests := []struct {
		name        string
		mockUsers   []*entities.UserOnlineStatus
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:      "successful get online users",
			mockUsers: []*entities.UserOnlineStatus{CreateTestUserOnlineStatus()},
			wantErr:   false,
		},
		{
			name:        "repository error",
			mockError:   errors.New("database error"),
			wantErr:     true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserOnlineStatusRepository(ctrl)
			mockRepo.EXPECT().GetOnlineUsers(gomock.Any()).Return(tt.mockUsers, tt.mockError)

			useCase := NewUserOnlineStatusUseCase(mockRepo)
			users, err := useCase.GetOnlineUsers(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, users)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, users)
				assert.Len(t, users, len(tt.mockUsers))
			}
		})
	}
}

func TestUserOnlineStatusUseCaseImpl_IsUserOnline(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockStatus  *entities.UserOnlineStatus
		mockError   error
		wantErr     bool
		expectedErr string
		expected    bool
	}{
		{
			name:       "user is online",
			userID:     "550e8400-e29b-41d4-a716-446655440000",
			mockStatus: CreateTestUserOnlineStatus(),
			wantErr:    false,
			expected:   true,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			mockError:   errors.New("invalid UUID length"),
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "status not found",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockError:   errors.New("status not found"),
			wantErr:     true,
			expectedErr: "status not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserOnlineStatusRepository(ctrl)
			mockRepo.EXPECT().GetByUserID(gomock.Any(), gomock.Any()).Return(tt.mockStatus, tt.mockError)

			useCase := NewUserOnlineStatusUseCase(mockRepo)
			isOnline, err := useCase.IsUserOnline(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, isOnline)
			}
		})
	}
}

func TestUserOnlineStatusUseCaseImpl_GetOnlineCount(t *testing.T) {
	tests := []struct {
		name        string
		mockCount   int64
		mockError   error
		wantErr     bool
		expectedErr string
		expected    int64
	}{
		{
			name:      "successful get online count",
			mockCount: 5,
			wantErr:   false,
			expected:  5,
		},
		{
			name:        "repository error",
			mockError:   errors.New("database error"),
			wantErr:     true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserOnlineStatusRepository(ctrl)
			mockRepo.EXPECT().GetOnlineCount(gomock.Any()).Return(tt.mockCount, tt.mockError)

			useCase := NewUserOnlineStatusUseCase(mockRepo)
			count, err := useCase.GetOnlineCount(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, count)
			}
		})
	}
}
