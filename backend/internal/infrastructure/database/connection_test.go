package database

import (
	"fmt"
	"testing"

	"github.com/atdevten/peace/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabaseManager(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		wantErr     bool
		expectedErr string
	}{
		{
			name: "create database manager with valid config",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						Host:     "localhost",
						Port:     "5432",
						User:     "postgres",
						Password: "password",
						DBName:   "peace",
						SSLMode:  "disable",
					},
				},
			},
			wantErr:     true, // Will fail due to no database connection
			expectedErr: "failed to connect to PostgreSQL",
		},
		{
			name: "create database manager with invalid host",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						Host:     "invalid_host",
						Port:     "5432",
						User:     "postgres",
						Password: "password",
						DBName:   "peace",
						SSLMode:  "disable",
					},
				},
			},
			wantErr:     true,
			expectedErr: "failed to connect to PostgreSQL",
		},
		{
			name: "create database manager with invalid port",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						Host:     "localhost",
						Port:     "9999",
						User:     "postgres",
						Password: "password",
						DBName:   "peace",
						SSLMode:  "disable",
					},
				},
			},
			wantErr:     true,
			expectedErr: "failed to connect to PostgreSQL",
		},
		{
			name: "create database manager with invalid user",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						Host:     "localhost",
						Port:     "5432",
						User:     "invalid_user",
						Password: "password",
						DBName:   "peace",
						SSLMode:  "disable",
					},
				},
			},
			wantErr:     true,
			expectedErr: "failed to connect to PostgreSQL",
		},
		{
			name: "create database manager with invalid password",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						Host:     "localhost",
						Port:     "5432",
						User:     "postgres",
						Password: "invalid_password",
						DBName:   "peace",
						SSLMode:  "disable",
					},
				},
			},
			wantErr:     true,
			expectedErr: "failed to connect to PostgreSQL",
		},
		{
			name: "create database manager with invalid database name",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						Host:     "localhost",
						Port:     "5432",
						User:     "postgres",
						Password: "password",
						DBName:   "invalid_db",
						SSLMode:  "disable",
					},
				},
			},
			wantErr:     true,
			expectedErr: "failed to connect to PostgreSQL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbManager, err := NewDatabaseManager(tt.config)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, dbManager)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, dbManager)
				assert.NotNil(t, dbManager.Postgres)

				// Test that we can ping the database
				sqlDB, err := dbManager.Postgres.DB()
				require.NoError(t, err)
				err = sqlDB.Ping()
				require.NoError(t, err)

				// Close the connection
				dbManager.Close()
			}
		})
	}
}

func TestNewDatabaseManager_NilConfig(t *testing.T) {
	// Test that nil config causes panic
	var panicOccurred bool
	var panicValue interface{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicOccurred = true
				panicValue = r
			}
		}()

		// This should panic due to nil pointer dereference
		NewDatabaseManager(nil)
	}()

	// Verify that panic occurred
	assert.True(t, panicOccurred, "Expected panic when passing nil config")
	assert.NotNil(t, panicValue, "Expected panic value to be non-nil")

	// Check that the panic is due to nil pointer dereference
	errStr := fmt.Sprintf("%v", panicValue)
	assert.Contains(t, errStr, "runtime error: invalid memory address or nil pointer dereference")
}

func TestNewDatabaseManager_EmptyDSN(t *testing.T) {
	config := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     "",
				Port:     "",
				User:     "",
				Password: "",
				DBName:   "",
				SSLMode:  "",
			},
		},
	}

	dbManager, err := NewDatabaseManager(config)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to PostgreSQL")
	assert.Nil(t, dbManager)
}

func TestDatabaseManager_GetMigrationVersion(t *testing.T) {
	// This test requires a real database connection
	// Skip if no test database is available
	t.Skip("Skipping migration test - requires real database")

	config := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "password",
				DBName:   "peace_test",
				SSLMode:  "disable",
			},
		},
	}

	dbManager, err := NewDatabaseManager(config)
	require.NoError(t, err)
	defer dbManager.Close()

	version, err := dbManager.GetMigrationVersion()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, version, int64(0))
}

func TestDatabaseManager_GetMigrationStatus(t *testing.T) {
	// This test requires a real database connection
	// Skip if no test database is available
	t.Skip("Skipping migration test - requires real database")

	config := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "password",
				DBName:   "peace_test",
				SSLMode:  "disable",
			},
		},
	}

	dbManager, err := NewDatabaseManager(config)
	require.NoError(t, err)
	defer dbManager.Close()

	status, err := dbManager.GetMigrationStatus()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, status, int64(0))
}

func TestConfig_GetPostgresDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected string
	}{
		{
			name: "generate DSN with all parameters",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						Host:     "localhost",
						Port:     "5432",
						User:     "postgres",
						Password: "password",
						DBName:   "peace",
						SSLMode:  "disable",
					},
				},
			},
			expected: "host=localhost port=5432 user=postgres password=password dbname=peace sslmode=disable",
		},
		{
			name: "generate DSN with minimal parameters",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						Host:     "localhost",
						Port:     "5432",
						User:     "postgres",
						Password: "password",
						DBName:   "peace",
						SSLMode:  "require",
					},
				},
			},
			expected: "host=localhost port=5432 user=postgres password=password dbname=peace sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.config.GetPostgresDSN()
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

func TestConfig_GetPostgresDSN_NilConfig(t *testing.T) {
	// Test that nil config causes panic
	var panicOccurred bool
	var panicValue interface{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicOccurred = true
				panicValue = r
			}
		}()

		// This should panic due to nil pointer dereference
		var config *config.Config
		config.GetPostgresDSN()
	}()

	// Verify that panic occurred
	assert.True(t, panicOccurred, "Expected panic when calling GetPostgresDSN on nil config")
	assert.NotNil(t, panicValue, "Expected panic value to be non-nil")

	// Check that the panic is due to nil pointer dereference
	errStr := fmt.Sprintf("%v", panicValue)
	assert.Contains(t, errStr, "runtime error: invalid memory address or nil pointer dereference")
}

func TestConfig_GetPostgresDSN_EmptyConfig(t *testing.T) {
	config := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{},
		},
	}

	dsn := config.GetPostgresDSN()
	expected := "host= port= user= password= dbname= sslmode="
	assert.Equal(t, expected, dsn)
}

func TestDatabaseManager_Integration(t *testing.T) {
	// This is an integration test that requires a real database
	// Skip if no test database is available
	t.Skip("Skipping integration test - requires real database")

	config := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "password",
				DBName:   "peace_test",
				SSLMode:  "disable",
			},
		},
	}

	// Test connection
	dbManager, err := NewDatabaseManager(config)
	require.NoError(t, err)
	defer dbManager.Close()

	// Test ping
	sqlDB, err := dbManager.Postgres.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	require.NoError(t, err)

	// Test query
	var result int
	err = dbManager.Postgres.Raw("SELECT 1").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 1, result)

	// Test migration version
	version, err := dbManager.GetMigrationVersion()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, version, int64(0))
}
