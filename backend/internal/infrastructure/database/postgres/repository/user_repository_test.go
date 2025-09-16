package repository

import (
	"context"
	"testing"

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

// setupUserTestDB creates an in-memory SQLite database for user testing
func setupUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Auto-migrate the models
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

func TestPostgreSQLUserRepository_Create(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewPostgreSQLUserRepository(db)

	tests := []struct {
		name    string
		user    *entities.User
		wantErr bool
	}{
		{
			name:    "successful creation",
			user:    helpers.CreateTestUser(),
			wantErr: false,
		},
		{
			name:    "successful Google user creation",
			user:    helpers.CreateTestGoogleUser(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.user)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify user was created
				var model models.User
				err = db.Where("id = ?", tt.user.ID().String()).First(&model).Error
				require.NoError(t, err)
				assert.Equal(t, tt.user.Email().String(), model.Email)
				assert.Equal(t, tt.user.Username().String(), model.Username)
			}
		})
	}
}

func TestPostgreSQLUserRepository_GetByID(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewPostgreSQLUserRepository(db)

	// Create a test user first
	testUser := helpers.CreateTestUser()
	err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name        string
		id          string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful get by ID",
			id:      testUser.ID().String(),
			wantErr: false,
		},
		{
			name:        "user not found",
			id:          "550e8400-e29b-41d4-a716-446655440001",
			wantErr:     true,
			expectedErr: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := value_objects.NewUserIDFromString(tt.id)
			require.NoError(t, err)

			user, err := repo.GetByID(context.Background(), userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, testUser.ID().String(), user.ID().String())
			}
		})
	}
}

func TestPostgreSQLUserRepository_GetByFilter(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewPostgreSQLUserRepository(db)

	// Create a test user first
	testUser := helpers.CreateTestUser()
	err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		filter  *repositories.UserFilter
		wantErr bool
	}{
		{
			name: "filter by email",
			filter: repositories.NewUserFilter(
				nil,
				testUser.Email(),
				nil,
			),
			wantErr: false,
		},
		{
			name: "filter by username",
			filter: repositories.NewUserFilter(
				nil,
				nil,
				testUser.Username(),
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByFilter(context.Background(), tt.filter)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, testUser.ID().String(), user.ID().String())
			}
		})
	}
}

func TestPostgreSQLUserRepository_Update(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewPostgreSQLUserRepository(db)

	// Create a test user first
	testUser := helpers.CreateTestUser()
	err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		user    *entities.User
		wantErr bool
	}{
		{
			name:    "successful update",
			user:    testUser,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(context.Background(), tt.user)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify user was updated
				var model models.User
				err = db.Where("id = ?", tt.user.ID().String()).First(&model).Error
				require.NoError(t, err)
				assert.Equal(t, tt.user.Email().String(), model.Email)
			}
		})
	}
}

func TestPostgreSQLUserRepository_Delete(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewPostgreSQLUserRepository(db)

	// Create a test user first
	testUser := helpers.CreateTestUser()
	err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name        string
		id          string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful deletion",
			id:      testUser.ID().String(),
			wantErr: false,
		},
		{
			name:    "user not found",
			id:      "550e8400-e29b-41d4-a716-446655440001",
			wantErr: false, // GORM Delete doesn't return error for non-existent records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := value_objects.NewUserIDFromString(tt.id)
			require.NoError(t, err)

			err = repo.Delete(context.Background(), userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)

				// Verify user was deleted
				var model models.User
				err = db.Where("id = ?", tt.id).First(&model).Error
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
			}
		})
	}
}

func TestPostgreSQLUserRepository_EmailExists(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewPostgreSQLUserRepository(db)

	// Create a test user first
	testUser := helpers.CreateTestUser()
	err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name     string
		email    *value_objects.Email
		wantErr  bool
		expected bool
	}{
		{
			name:     "email exists",
			email:    testUser.Email(),
			wantErr:  false,
			expected: true,
		},
		{
			name:     "email does not exist",
			email:    helpers.CreateTestEmail("nonexistent@example.com"),
			wantErr:  false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.EmailExists(context.Background(), *tt.email)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			}
		})
	}
}
