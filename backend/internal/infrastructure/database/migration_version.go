package database

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// MigrationInfo represents information about a migration
type MigrationInfo struct {
	Version     int64
	Name        string
	Description string
	CreatedAt   time.Time
	AppliedAt   *time.Time
	IsApplied   bool
}

// MigrationVersionManager manages migration versions
type MigrationVersionManager struct {
	db *gorm.DB
}

// NewMigrationVersionManager creates a new migration version manager
func NewMigrationVersionManager(db *gorm.DB) *MigrationVersionManager {
	return &MigrationVersionManager{db: db}
}

// GetMigrationInfo returns detailed information about all migrations
func (mvm *MigrationVersionManager) GetMigrationInfo() ([]MigrationInfo, error) {
	// Get current database version
	currentVersion, err := mvm.getCurrentVersion()
	if err != nil {
		return nil, err
	}

	// Get all migration files
	migrationFiles, err := mvm.getMigrationFiles()
	if err != nil {
		return nil, err
	}

	var migrations []MigrationInfo
	for _, file := range migrationFiles {
		version := mvm.extractVersionFromFilename(file.Name())
		description := mvm.extractDescriptionFromFilename(file.Name())

		migration := MigrationInfo{
			Version:     version,
			Name:        file.Name(),
			Description: description,
			CreatedAt:   time.Now(), // Using current time as fallback
			IsApplied:   version <= currentVersion,
		}

		// If applied, get applied time
		if migration.IsApplied {
			appliedAt, err := mvm.getMigrationAppliedTime(version)
			if err == nil {
				migration.AppliedAt = &appliedAt
			}
		}

		migrations = append(migrations, migration)
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// GetCurrentVersion returns the current migration version
func (mvm *MigrationVersionManager) GetCurrentVersion() (int64, error) {
	return mvm.getCurrentVersion()
}

// GetPendingMigrations returns migrations that haven't been applied
func (mvm *MigrationVersionManager) GetPendingMigrations() ([]MigrationInfo, error) {
	allMigrations, err := mvm.GetMigrationInfo()
	if err != nil {
		return nil, err
	}

	var pending []MigrationInfo
	for _, migration := range allMigrations {
		if !migration.IsApplied {
			pending = append(pending, migration)
		}
	}

	return pending, nil
}

// GetAppliedMigrations returns migrations that have been applied
func (mvm *MigrationVersionManager) GetAppliedMigrations() ([]MigrationInfo, error) {
	allMigrations, err := mvm.GetMigrationInfo()
	if err != nil {
		return nil, err
	}

	var applied []MigrationInfo
	for _, migration := range allMigrations {
		if migration.IsApplied {
			applied = append(applied, migration)
		}
	}

	return applied, nil
}

// PrintMigrationStatus prints a formatted migration status
func (mvm *MigrationVersionManager) PrintMigrationStatus() error {
	migrations, err := mvm.GetMigrationInfo()
	if err != nil {
		return err
	}

	currentVersion, err := mvm.getCurrentVersion()
	if err != nil {
		return err
	}

	fmt.Printf("Database Migration Status\n")
	fmt.Printf("========================\n")
	fmt.Printf("Current Version: %d\n\n", currentVersion)

	for _, migration := range migrations {
		status := "❌"
		if migration.IsApplied {
			status = "✅"
		}

		fmt.Printf("%s %03d - %s\n", status, migration.Version, migration.Description)

		if migration.IsApplied && migration.AppliedAt != nil {
			fmt.Printf("    Applied: %s\n", migration.AppliedAt.Format("2006-01-02 15:04:05"))
		}
	}

	return nil
}

// getCurrentVersion gets the current migration version from database
func (mvm *MigrationVersionManager) getCurrentVersion() (int64, error) {
	sqlDB, err := mvm.db.DB()
	if err != nil {
		return 0, err
	}

	return goose.GetDBVersion(sqlDB)
}

// getMigrationFiles gets all migration files from the migrations directory
func (mvm *MigrationVersionManager) getMigrationFiles() ([]os.DirEntry, error) {
	migrationsDir := "migrations"

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return []os.DirEntry{}, nil
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []os.DirEntry
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file)
		}
	}

	return migrationFiles, nil
}

// extractVersionFromFilename extracts version number from migration filename
func (mvm *MigrationVersionManager) extractVersionFromFilename(filename string) int64 {
	parts := strings.Split(filename, "_")
	if len(parts) > 0 {
		if version, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
			return version
		}
	}
	return 0
}

// extractDescriptionFromFilename extracts description from migration filename
func (mvm *MigrationVersionManager) extractDescriptionFromFilename(filename string) string {
	// Remove .sql extension
	name := strings.TrimSuffix(filename, ".sql")

	// Remove version prefix (e.g., "001_")
	parts := strings.SplitN(name, "_", 2)
	if len(parts) > 1 {
		return strings.ReplaceAll(parts[1], "_", " ")
	}

	return name
}

// getMigrationAppliedTime gets when a migration was applied
func (mvm *MigrationVersionManager) getMigrationAppliedTime(version int64) (time.Time, error) {
	sqlDB, err := mvm.db.DB()
	if err != nil {
		return time.Time{}, err
	}

	var appliedAt time.Time
	query := `SELECT executed_at FROM goose_db_version WHERE version_id = $1`
	err = sqlDB.QueryRow(query, version).Scan(&appliedAt)

	return appliedAt, err
}
