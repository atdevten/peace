package repository

import (
	"context"
	"testing"
	"time"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
	"github.com/atdevten/peace/internal/infrastructure/database/postgres/models"
	"github.com/atdevten/peace/testutils/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Auto-migrate the models
	err = db.AutoMigrate(&models.MentalHealthRecord{})
	require.NoError(t, err)

	return db
}

func TestPostgreSQLMentalHealthRecordRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgreSQLMentalHealthRecordRepository(db)

	tests := []struct {
		name    string
		record  *entities.MentalHealthRecord
		wantErr bool
	}{
		{
			name:    "successful creation",
			record:  helpers.CreateTestMentalHealthRecord(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.record)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify record was created
				var model models.MentalHealthRecord
				err = db.Where("id = ?", tt.record.ID().String()).First(&model).Error
				require.NoError(t, err)
				assert.Equal(t, tt.record.UserID().String(), model.UserID)
				assert.Equal(t, tt.record.HappyLevel().Value(), model.HappyLevel)
				assert.Equal(t, tt.record.EnergyLevel().Value(), model.EnergyLevel)
			}
		})
	}
}

func TestPostgreSQLMentalHealthRecordRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgreSQLMentalHealthRecordRepository(db)

	// Create a test record first
	testRecord := helpers.CreateTestMentalHealthRecord()
	err := repo.Create(context.Background(), testRecord)
	require.NoError(t, err)

	tests := []struct {
		name        string
		id          string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful get by ID",
			id:      testRecord.ID().String(),
			wantErr: false,
		},
		{
			name:        "record not found",
			id:          "550e8400-e29b-41d4-a716-446655440001",
			wantErr:     true,
			expectedErr: "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recordID, err := value_objects.NewMentalHealthRecordIDFromString(tt.id)
			require.NoError(t, err)

			record, err := repo.GetByID(context.Background(), recordID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, record)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, record)
				assert.Equal(t, testRecord.ID().String(), record.ID().String())
			}
		})
	}
}

func TestPostgreSQLMentalHealthRecordRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgreSQLMentalHealthRecordRepository(db)

	// Create a test record first
	testRecord := helpers.CreateTestMentalHealthRecord()
	err := repo.Create(context.Background(), testRecord)
	require.NoError(t, err)

	tests := []struct {
		name    string
		record  *entities.MentalHealthRecord
		wantErr bool
	}{
		{
			name:    "successful update",
			record:  testRecord,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(context.Background(), tt.record)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify record was updated
				var model models.MentalHealthRecord
				err = db.Where("id = ?", tt.record.ID().String()).First(&model).Error
				require.NoError(t, err)
				assert.Equal(t, tt.record.UserID().String(), model.UserID)
			}
		})
	}
}

func TestPostgreSQLMentalHealthRecordRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgreSQLMentalHealthRecordRepository(db)

	// Create a test record first
	testRecord := helpers.CreateTestMentalHealthRecord()
	err := repo.Create(context.Background(), testRecord)
	require.NoError(t, err)

	tests := []struct {
		name        string
		id          string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful deletion",
			id:      testRecord.ID().String(),
			wantErr: false,
		},
		{
			name:        "record not found",
			id:          "550e8400-e29b-41d4-a716-446655440001",
			wantErr:     true,
			expectedErr: "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recordID, err := value_objects.NewMentalHealthRecordIDFromString(tt.id)
			require.NoError(t, err)

			err = repo.Delete(context.Background(), recordID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)

				// Verify record was deleted
				var model models.MentalHealthRecord
				err = db.Where("id = ?", tt.id).First(&model).Error
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
			}
		})
	}
}

func TestPostgreSQLMentalHealthRecordRepository_GetByFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgreSQLMentalHealthRecordRepository(db)

	// Create test records
	userID := helpers.CreateTestUserID()
	record1 := helpers.CreateTestMentalHealthRecordWithUserID(userID)
	record2 := helpers.CreateTestMentalHealthRecordWithUserID(userID)

	err := repo.Create(context.Background(), record1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), record2)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filter   *repositories.MentalHealthRecordFilter
		wantErr  bool
		minCount int
	}{
		{
			name: "filter by user ID",
			filter: &repositories.MentalHealthRecordFilter{
				UserID: userID,
			},
			wantErr:  false,
			minCount: 1,
		},
		{
			name: "filter by date range",
			filter: &repositories.MentalHealthRecordFilter{
				UserID:    userID,
				StartedAt: &[]time.Time{time.Now().AddDate(0, 0, -1)}[0],
				EndedAt:   &[]time.Time{time.Now().AddDate(0, 0, 1)}[0],
			},
			wantErr:  false,
			minCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := repo.GetByFilter(context.Background(), tt.filter)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(records), tt.minCount)
			}
		})
	}
}

func TestPostgreSQLMentalHealthRecordRepository_GetDistinctDatesForUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgreSQLMentalHealthRecordRepository(db)

	// Create test records
	userID := helpers.CreateTestUserID()
	record1 := helpers.CreateTestMentalHealthRecordWithUserID(userID)
	record2 := helpers.CreateTestMentalHealthRecordWithUserID(userID)

	err := repo.Create(context.Background(), record1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), record2)
	require.NoError(t, err)

	tests := []struct {
		name     string
		userID   *value_objects.UserID
		wantErr  bool
		minCount int
	}{
		{
			name:     "successful get distinct dates",
			userID:   userID,
			wantErr:  false,
			minCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dates, err := repo.GetDistinctDatesForUser(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(dates), tt.minCount)
			}
		})
	}
}
