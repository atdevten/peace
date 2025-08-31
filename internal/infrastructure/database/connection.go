package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atdevten/peace/internal/infrastructure/config"

	"github.com/pressly/goose/v3"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseManager struct {
	Postgres *gorm.DB
}

func NewDatabaseManager(cfg *config.Config) (*DatabaseManager, error) {

	// Initialize PostgreSQL
	postgresDB, err := gorm.Open(gormPostgres.Open(cfg.GetPostgresDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := postgresDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.Postgres.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.Postgres.ConnMaxLifetime)

	// Run migrations using Goose with versioning
	if err := runGooseMigrations(postgresDB); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DatabaseManager{
		Postgres: postgresDB,
	}, nil
}

// runGooseMigrations runs database migrations using Goose with versioning
func runGooseMigrations(db *gorm.DB) error {
	// Set the dialect for Goose
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Get the underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Get the executable path to find the project root
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory of the executable
	execDir := filepath.Dir(execPath)

	// Find the project root by looking for migrations directory
	// Start from executable directory and go up until we find migrations
	var migrationsPath string
	currentDir := execDir

	for i := 0; i < 10; i++ { // Limit to prevent infinite loop
		testPath := filepath.Join(currentDir, "migrations")
		if _, err := os.Stat(testPath); err == nil {
			migrationsPath = testPath
			break
		}

		// Go up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// We've reached the root directory
			break
		}
		currentDir = parentDir
	}

	if migrationsPath == "" {
		// Fallback: try relative to current working directory
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		testPath := filepath.Join(wd, "migrations")
		if _, err := os.Stat(testPath); err == nil {
			migrationsPath = testPath
		} else {
			return fmt.Errorf("migrations directory not found. Tried paths relative to executable: %s and working directory: %s", execDir, wd)
		}
	}

	// Run migrations
	if err = goose.Up(sqlDB, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// GetMigrationVersion returns the current migration version
func (dm *DatabaseManager) GetMigrationVersion() (int64, error) {
	sqlDB, err := dm.Postgres.DB()
	if err != nil {
		return 0, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	version, err := goose.GetDBVersion(sqlDB)
	if err != nil {
		return 0, err
	}

	return version, nil
}

// GetMigrationVersionManager returns a migration version manager
func (dm *DatabaseManager) GetMigrationVersionManager() *MigrationVersionManager {
	return NewMigrationVersionManager(dm.Postgres)
}

// GetMigrationStatus returns detailed migration status
func (dm *DatabaseManager) GetMigrationStatus() (int64, error) {
	sqlDB, err := dm.Postgres.DB()
	if err != nil {
		return 0, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	version, err := goose.GetDBVersion(sqlDB)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func (dm *DatabaseManager) Close() {
	if dm.Postgres != nil {
		sqlDB, err := dm.Postgres.DB()
		if err != nil {
			panic(err)
		}
		if err = sqlDB.Close(); err != nil {
			panic(err)
		}
	}
}
