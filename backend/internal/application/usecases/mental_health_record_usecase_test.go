package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/atdevten/peace/internal/application/commands"
	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/testutils/helpers"
	repositories "github.com/atdevten/peace/testutils/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMentalHealthRecordUseCaseImpl_Create(t *testing.T) {
	tests := []struct {
		name        string
		command     commands.CreateMentalHealthRecordCommand
		mockRecord  *entities.MentalHealthRecord
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful creation",
			command: commands.CreateMentalHealthRecordCommand{
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				HappyLevel:  5,
				EnergyLevel: 7,
				Notes:       helpers.StringPtr("Feeling good today"),
				Status:      "public",
			},
			mockRecord: helpers.CreateTestMentalHealthRecord(),
			wantErr:    false,
		},
		{
			name: "invalid happy level",
			command: commands.CreateMentalHealthRecordCommand{
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				HappyLevel:  15, // Invalid level
				EnergyLevel: 7,
				Notes:       helpers.StringPtr("Feeling good today"),
				Status:      "public",
			},
			wantErr:     true,
			expectedErr: "happy level must be between 1 and 10",
		},
		{
			name: "invalid energy level",
			command: commands.CreateMentalHealthRecordCommand{
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				HappyLevel:  5,
				EnergyLevel: 15, // Invalid level
				Notes:       helpers.StringPtr("Feeling good today"),
				Status:      "public",
			},
			wantErr:     true,
			expectedErr: "energy level must be between 1 and 10",
		},
		{
			name: "repository error",
			command: commands.CreateMentalHealthRecordCommand{
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				HappyLevel:  5,
				EnergyLevel: 7,
				Notes:       helpers.StringPtr("Feeling good today"),
				Status:      "public",
			},
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
			mockRepo := repositories.NewMockMentalHealthRecordRepository(ctrl)
			if tt.command.HappyLevel >= 1 && tt.command.HappyLevel <= 10 &&
				tt.command.EnergyLevel >= 1 && tt.command.EnergyLevel <= 10 &&
				tt.command.Status == "public" {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.mockError)
				if tt.mockError == nil {
					mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockRecord, nil)
				}
			}

			useCase := NewMentalHealthRecordUseCase(mockRepo)
			record, err := useCase.Create(context.Background(), tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, record)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, record)
			}
		})
	}
}

func TestMentalHealthRecordUseCaseImpl_Update(t *testing.T) {
	tests := []struct {
		name        string
		command     commands.UpdateMentalHealthRecordCommand
		mockRecord  *entities.MentalHealthRecord
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful update",
			command: commands.UpdateMentalHealthRecordCommand{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				HappyLevel:  8,
				EnergyLevel: 6,
				Notes:       helpers.StringPtr("Updated feeling"),
				Status:      "public",
			},
			mockRecord: helpers.CreateTestMentalHealthRecord(),
			wantErr:    false,
		},
		{
			name: "invalid record ID",
			command: commands.UpdateMentalHealthRecordCommand{
				ID:          "invalid-id",
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				HappyLevel:  8,
				EnergyLevel: 6,
				Notes:       helpers.StringPtr("Updated feeling"),
				Status:      "public",
			},
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name: "record not found",
			command: commands.UpdateMentalHealthRecordCommand{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				HappyLevel:  8,
				EnergyLevel: 6,
				Notes:       helpers.StringPtr("Updated feeling"),
				Status:      "public",
			},
			mockError:   errors.New("record not found"),
			wantErr:     true,
			expectedErr: "record not found",
		},
		{
			name: "unauthorized - different user",
			command: commands.UpdateMentalHealthRecordCommand{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
				UserID:      "different-user-id",
				HappyLevel:  8,
				EnergyLevel: 6,
				Notes:       helpers.StringPtr("Updated feeling"),
				Status:      "public",
			},
			mockRecord:  helpers.CreateTestMentalHealthRecord(),
			wantErr:     true,
			expectedErr: "unauthorized: user does not own this record",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockMentalHealthRecordRepository(ctrl)
			if tt.command.ID != "invalid-id" {
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockRecord, tt.mockError)
				if tt.mockError == nil && tt.command.UserID == "550e8400-e29b-41d4-a716-446655440000" {
					mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			useCase := NewMentalHealthRecordUseCase(mockRepo)
			record, err := useCase.Update(context.Background(), tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, record)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, record)
			}
		})
	}
}

func TestMentalHealthRecordUseCaseImpl_Delete(t *testing.T) {
	tests := []struct {
		name        string
		command     commands.DeleteMentalHealthRecordCommand
		mockRecord  *entities.MentalHealthRecord
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful deletion",
			command: commands.DeleteMentalHealthRecordCommand{
				ID:     "550e8400-e29b-41d4-a716-446655440000",
				UserID: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockRecord: helpers.CreateTestMentalHealthRecord(),
			wantErr:    false,
		},
		{
			name: "invalid record ID",
			command: commands.DeleteMentalHealthRecordCommand{
				ID:     "invalid-id",
				UserID: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name: "record not found",
			command: commands.DeleteMentalHealthRecordCommand{
				ID:     "550e8400-e29b-41d4-a716-446655440000",
				UserID: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockError:   errors.New("record not found"),
			wantErr:     true,
			expectedErr: "record not found",
		},
		{
			name: "unauthorized - different user",
			command: commands.DeleteMentalHealthRecordCommand{
				ID:     "550e8400-e29b-41d4-a716-446655440000",
				UserID: "different-user-id",
			},
			mockRecord:  helpers.CreateTestMentalHealthRecord(),
			wantErr:     true,
			expectedErr: "unauthorized: user does not own this record",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockMentalHealthRecordRepository(ctrl)
			if tt.command.ID != "invalid-id" {
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockRecord, tt.mockError)
				if tt.mockError == nil && tt.command.UserID == "550e8400-e29b-41d4-a716-446655440000" {
					mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			useCase := NewMentalHealthRecordUseCase(mockRepo)
			err := useCase.Delete(context.Background(), tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMentalHealthRecordUseCaseImpl_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		userID      string
		mockRecord  *entities.MentalHealthRecord
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:       "successful get",
			id:         "550e8400-e29b-41d4-a716-446655440000",
			userID:     "550e8400-e29b-41d4-a716-446655440000",
			mockRecord: helpers.CreateTestMentalHealthRecord(),
			wantErr:    false,
		},
		{
			name:        "invalid record ID",
			id:          "invalid-id",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "record not found",
			id:          "550e8400-e29b-41d4-a716-446655440000",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockError:   errors.New("record not found"),
			wantErr:     true,
			expectedErr: "record not found",
		},
		{
			name:        "unauthorized - different user",
			id:          "550e8400-e29b-41d4-a716-446655440000",
			userID:      "different-user-id",
			mockRecord:  helpers.CreateTestMentalHealthRecord(),
			wantErr:     true,
			expectedErr: "unauthorized: user does not own this record",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockMentalHealthRecordRepository(ctrl)
			if tt.id != "invalid-id" {
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockRecord, tt.mockError)
			}

			useCase := NewMentalHealthRecordUseCase(mockRepo)
			record, err := useCase.GetByID(context.Background(), tt.id, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, record)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, record)
			}
		})
	}
}

func TestMentalHealthRecordUseCaseImpl_GetByCondition(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		startedAt   *string
		endedAt     *string
		mockRecords []*entities.MentalHealthRecord
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "successful get with no date range",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			mockRecords: []*entities.MentalHealthRecord{helpers.CreateTestMentalHealthRecord()},
			wantErr:     false,
		},
		{
			name:        "successful get with date range",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			startedAt:   helpers.StringPtr("2023-01-01T00:00:00Z"),
			endedAt:     helpers.StringPtr("2023-12-31T23:59:59Z"),
			mockRecords: []*entities.MentalHealthRecord{helpers.CreateTestMentalHealthRecord()},
			wantErr:     false,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid-id",
			wantErr:     true,
			expectedErr: "invalid UUID length",
		},
		{
			name:        "invalid date format",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			startedAt:   helpers.StringPtr("invalid-date"),
			wantErr:     true,
			expectedErr: "invalid start date format",
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
			mockRepo := repositories.NewMockMentalHealthRecordRepository(ctrl)
			if tt.userID != "invalid-id" && (tt.startedAt == nil || *tt.startedAt != "invalid-date") {
				mockRepo.EXPECT().GetByFilter(gomock.Any(), gomock.Any()).Return(tt.mockRecords, tt.mockError)
			}

			useCase := NewMentalHealthRecordUseCase(mockRepo)
			records, err := useCase.GetByCondition(context.Background(), tt.userID, tt.startedAt, tt.endedAt)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, records)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, records)
				assert.Len(t, records, len(tt.mockRecords))
			}
		})
	}
}
