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

func TestUserUseCaseImpl_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockUser    *entities.User
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "successful get user",
			userID:   "550e8400-e29b-41d4-a716-446655440000",
			mockUser: helpers.CreateTestUser(),
			wantErr:  false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "user not found",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockError:   errors.New("user not found"),
			wantErr:     true,
			expectedErr: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)
			if tt.userID != "invalid-id" {
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockUser, tt.mockError)
			}

			useCase := NewUserUseCase(mockRepo)
			user, err := useCase.GetByID(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.mockUser.Email().String(), user.Email().String())
			}
		})
	}
}

func TestUserUseCaseImpl_UpdateProfile(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		firstName   *string
		lastName    *string
		mockUser    *entities.User
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:      "successful profile update",
			userID:    "550e8400-e29b-41d4-a716-446655440000",
			firstName: helpers.StringPtr("Updated"),
			lastName:  helpers.StringPtr("Name"),
			mockUser:  helpers.CreateTestUser(),
			wantErr:   false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			firstName:   helpers.StringPtr("Updated"),
			lastName:    helpers.StringPtr("Name"),
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "user not found",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			firstName:   helpers.StringPtr("Updated"),
			lastName:    helpers.StringPtr("Name"),
			mockError:   errors.New("user not found"),
			wantErr:     true,
			expectedErr: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)
			if tt.userID != "invalid-id" {
				// Mock GetByID call
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockUser, tt.mockError)
				// Mock Update call if no error in GetByID
				if tt.mockError == nil {
					mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			useCase := NewUserUseCase(mockRepo)
			user, err := useCase.UpdateProfile(context.Background(), tt.userID, tt.firstName, tt.lastName)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}

func TestUserUseCaseImpl_UpdatePassword(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		newPassword string
		mockUser    *entities.User
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "successful password update",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			newPassword: "NewPassword123",
			mockUser:    helpers.CreateTestUser(),
			wantErr:     false,
		},
		{
			name:        "empty password",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			newPassword: "",
			wantErr:     true,
			expectedErr: "new password is required",
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			newPassword: "NewPassword123",
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)
			if tt.userID != "invalid-id" && tt.newPassword != "" {
				// Mock GetByID call
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockUser, tt.mockError)
				// Mock Update call if no error in GetByID
				if tt.mockError == nil {
					mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			useCase := NewUserUseCase(mockRepo)
			err := useCase.UpdatePassword(context.Background(), tt.userID, tt.newPassword)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserUseCaseImpl_Deactivate(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockUser    *entities.User
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "successful deactivation",
			userID:   "550e8400-e29b-41d4-a716-446655440000",
			mockUser: helpers.CreateTestUser(),
			wantErr:  false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "user not found",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockError:   errors.New("user not found"),
			wantErr:     true,
			expectedErr: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)
			if tt.userID != "invalid-id" {
				// Mock GetByID call
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockUser, tt.mockError)
				// Mock Update call if no error in GetByID
				if tt.mockError == nil {
					mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			useCase := NewUserUseCase(mockRepo)
			err := useCase.Deactivate(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserUseCaseImpl_Delete(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful deletion",
			userID:  "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "user not found",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockError:   errors.New("user not found"),
			wantErr:     true,
			expectedErr: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockUserRepository(ctrl)
			if tt.userID != "invalid-id" {
				mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.mockError)
			}

			useCase := NewUserUseCase(mockRepo)
			err := useCase.Delete(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
