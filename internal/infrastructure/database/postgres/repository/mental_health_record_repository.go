package repository

import (
	"context"
	"fmt"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
	"github.com/atdevten/peace/internal/infrastructure/database/postgres/models"
	"gorm.io/gorm"
)

type PostgreSQLMentalHealthRecordRepository struct {
	db *gorm.DB
}

func NewPostgreSQLMentalHealthRecordRepository(db *gorm.DB) repositories.MentalHealthRecordRepository {
	return &PostgreSQLMentalHealthRecordRepository{
		db: db,
	}
}

func (r *PostgreSQLMentalHealthRecordRepository) Create(ctx context.Context, record *entities.MentalHealthRecord) error {
	model := models.MentalHealthRecord{
		ID:          record.ID().String(),
		UserID:      record.UserID().String(),
		HappyLevel:  record.HappyLevel().Value(),
		EnergyLevel: record.EnergyLevel().Value(),
		Notes:       record.Notes(),
		Status:      record.Status().String(),
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("r.db.Create: %w", err)
	}

	return nil
}

func (r *PostgreSQLMentalHealthRecordRepository) GetByID(ctx context.Context, id *value_objects.MentalHealthRecordID) (*entities.MentalHealthRecord, error) {
	var model models.MentalHealthRecord

	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("r.db.First: %w", err)
	}

	// Convert model to entity
	userID, err := value_objects.NewUserIDFromString(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewUserIDFromString: %w", err)
	}

	happyLevel, err := value_objects.NewHappyLevel(model.HappyLevel)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewHappyLevel: %w", err)
	}

	energyLevel, err := value_objects.NewEnergyLevel(model.EnergyLevel)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewEnergyLevel: %w", err)
	}

	recordID, err := value_objects.NewMentalHealthRecordIDFromString(model.ID)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewMentalHealthRecordIDFromString: %w", err)
	}

	status, err := value_objects.NewMentalHealthRecordStatus(model.Status)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewMentalHealthRecordStatus: %w", err)
	}

	record := entities.NewMentalHealthRecordFromExisting(
		recordID,
		userID,
		happyLevel,
		energyLevel,
		model.Notes,
		status,
		model.CreatedAt,
		model.UpdatedAt,
		model.DeletedAt,
	)

	return record, nil
}

func (r *PostgreSQLMentalHealthRecordRepository) Delete(ctx context.Context, id *value_objects.MentalHealthRecordID) error {
	result := r.db.WithContext(ctx).Delete(&models.MentalHealthRecord{}, id.String())
	if result.Error != nil {
		return fmt.Errorf("r.db.Delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (r *PostgreSQLMentalHealthRecordRepository) GetByFilter(ctx context.Context, filter *repositories.MentalHealthRecordFilter) ([]*entities.MentalHealthRecord, error) {
	var modelList []models.MentalHealthRecord

	// Build query with user_id filter
	query := r.db.WithContext(ctx).Model(&models.MentalHealthRecord{})
	if filter.UserID != nil {
		query = query.Where("user_id = ?", filter.UserID.String())
	}

	// Add time range filters if provided
	if filter.StartedAt != nil {
		query = query.Where("created_at >= ?", *filter.StartedAt)
	}
	if filter.EndedAt != nil {
		query = query.Where("created_at <= ?", *filter.EndedAt)
	}

	// Order by created_at; default DESC, allow override
	if filter.OrderDesc {
		query = query.Order("created_at ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	// Pagination
	if filter.Offset != nil {
		query = query.Offset(*filter.Offset)
	}
	if filter.Limit != nil {
		query = query.Limit(*filter.Limit)
	}

	fmt.Println("query", query.Statement.SQL.String())

	// Get all records matching the condition
	if err := query.Find(&modelList).Error; err != nil {
		return nil, fmt.Errorf("r.db.Find: %w", err)
	}

	// Convert models to entities
	var records []*entities.MentalHealthRecord
	for _, model := range modelList {
		userIDVO, err := value_objects.NewUserIDFromString(model.UserID)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewUserIDFromString: %w", err)
		}

		happyLevel, err := value_objects.NewHappyLevel(model.HappyLevel)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewHappyLevel: %w", err)
		}

		energyLevel, err := value_objects.NewEnergyLevel(model.EnergyLevel)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewEnergyLevel: %w", err)
		}

		recordID, err := value_objects.NewMentalHealthRecordIDFromString(model.ID)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewMentalHealthRecordIDFromString: %w", err)
		}

		status, err := value_objects.NewMentalHealthRecordStatus(model.Status)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewMentalHealthRecordStatus: %w", err)
		}

		record := entities.NewMentalHealthRecordFromExisting(
			recordID,
			userIDVO,
			happyLevel,
			energyLevel,
			model.Notes,
			status,
			model.CreatedAt,
			model.UpdatedAt,
			model.DeletedAt,
		)

		records = append(records, record)
	}

	return records, nil
}

func (r *PostgreSQLMentalHealthRecordRepository) GetAll(ctx context.Context) ([]*entities.MentalHealthRecord, error) {
	var modelList []models.MentalHealthRecord

	if err := r.db.WithContext(ctx).Find(&modelList).Error; err != nil {
		return nil, fmt.Errorf("r.db.Find: %w", err)
	}

	// Convert models to entities
	var records []*entities.MentalHealthRecord
	for _, model := range modelList {
		userIDVO, err := value_objects.NewUserIDFromString(model.UserID)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewUserIDFromString: %w", err)
		}

		happyLevel, err := value_objects.NewHappyLevel(model.HappyLevel)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewHappyLevel: %w", err)
		}

		energyLevel, err := value_objects.NewEnergyLevel(model.EnergyLevel)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewEnergyLevel: %w", err)
		}

		recordID, err := value_objects.NewMentalHealthRecordIDFromString(model.ID)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewMentalHealthRecordIDFromString: %w", err)
		}

		status, err := value_objects.NewMentalHealthRecordStatus(model.Status)
		if err != nil {
			return nil, fmt.Errorf("value_objects.NewMentalHealthRecordStatus: %w", err)
		}

		record := entities.NewMentalHealthRecordFromExisting(
			recordID,
			userIDVO,
			happyLevel,
			energyLevel,
			model.Notes,
			status,
			model.CreatedAt,
			model.UpdatedAt,
			model.DeletedAt,
		)

		records = append(records, record)
	}

	return records, nil
}

func (r *PostgreSQLMentalHealthRecordRepository) Update(ctx context.Context, record *entities.MentalHealthRecord) error {
	model := models.MentalHealthRecord{
		ID:          record.ID().String(),
		UserID:      record.UserID().String(),
		HappyLevel:  record.HappyLevel().Value(),
		EnergyLevel: record.EnergyLevel().Value(),
		Notes:       record.Notes(),
		Status:      record.Status().String(),
	}

	if err := r.db.WithContext(ctx).Model(&models.MentalHealthRecord{}).Where("id = ?", model.ID).Updates(&model).Error; err != nil {
		return fmt.Errorf("r.db.Updates: %w", err)
	}

	return nil
}

func (r *PostgreSQLMentalHealthRecordRepository) GetDistinctDatesForUser(ctx context.Context, userID *value_objects.UserID) ([]string, error) {
	var records []models.MentalHealthRecord

	// Use GORM to get records ordered by created_at DESC (leverages index)
	// Only select created_at to minimize data transfer
	err := r.db.WithContext(ctx).
		Select("created_at").
		Where("user_id = ? AND deleted_at IS NULL", userID.String()).
		Order("created_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("r.db.Find: %w", err)
	}

	// Extract unique dates in application layer (more efficient than DB function)
	seenDates := make(map[string]bool)
	var dates []string

	for _, record := range records {
		dateStr := record.CreatedAt.Format("2006-01-02")
		if !seenDates[dateStr] {
			seenDates[dateStr] = true
			dates = append(dates, dateStr)
		}
	}

	return dates, nil
}
