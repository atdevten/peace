package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/atdevten/peace/internal/application/commands"
	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
	"github.com/atdevten/peace/internal/pkg/timeutil"
)

type MentalHealthRecordUseCase interface {
	Create(ctx context.Context, command commands.CreateMentalHealthRecordCommand) (*entities.MentalHealthRecord, error)
	Update(ctx context.Context, command commands.UpdateMentalHealthRecordCommand) (*entities.MentalHealthRecord, error)
	Delete(ctx context.Context, command commands.DeleteMentalHealthRecordCommand) error
	GetByID(ctx context.Context, id string, userID string) (*entities.MentalHealthRecord, error)
	GetByCondition(ctx context.Context, userID string, startedAt *string, endedAt *string) ([]*entities.MentalHealthRecord, error)
	GetHeatmap(ctx context.Context, command *commands.GetMentalHealthHeatmapCommand) (*commands.MentalHealthHeatmapResult, error)
	GetStreak(ctx context.Context, command *commands.GetMentalHealthStreakCommand) (*commands.MentalHealthStreakResult, error)
}

type MentalHealthRecordUseCaseImpl struct {
	recordRepo repositories.MentalHealthRecordRepository
}

func NewMentalHealthRecordUseCase(recordRepo repositories.MentalHealthRecordRepository) MentalHealthRecordUseCase {
	return &MentalHealthRecordUseCaseImpl{
		recordRepo: recordRepo,
	}
}

func (uc *MentalHealthRecordUseCaseImpl) Create(ctx context.Context, command commands.CreateMentalHealthRecordCommand) (*entities.MentalHealthRecord, error) {
	// Create new record
	record, err := entities.NewMentalHealthRecord(
		command.UserID,
		command.HappyLevel,
		command.EnergyLevel,
		command.Notes,
		command.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("entities.NewMentalHealthRecord: %w", err)
	}

	// Create in repository
	if err := uc.recordRepo.Create(ctx, record); err != nil {
		return nil, fmt.Errorf("uc.recordRepo.Create: %w", err)
	}

	// Get record
	newRecord, err := uc.recordRepo.GetByID(ctx, record.ID())
	if err != nil {
		return nil, fmt.Errorf("uc.recordRepo.GetByID: %w", err)
	}

	return newRecord, nil
}

func (uc *MentalHealthRecordUseCaseImpl) Update(ctx context.Context, command commands.UpdateMentalHealthRecordCommand) (*entities.MentalHealthRecord, error) {
	// Parse record ID
	recordID, err := value_objects.NewMentalHealthRecordIDFromString(command.ID)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewMentalHealthRecordIDFromString: %w", err)
	}

	// Find existing record
	existingRecord, err := uc.recordRepo.GetByID(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("uc.recordRepo.GetByID: %w", err)
	}

	// Check if user owns this record
	if existingRecord.UserID().String() != command.UserID {
		return nil, fmt.Errorf("unauthorized: user does not own this record")
	}

	// Create new happy level
	happyLevel, err := value_objects.NewHappyLevel(command.HappyLevel)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewHappyLevel: %w", err)
	}

	// Create new energy level
	energyLevel, err := value_objects.NewEnergyLevel(command.EnergyLevel)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewEnergyLevel: %w", err)
	}

	status, err := value_objects.NewMentalHealthRecordStatus(command.Status)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewMentalHealthRecordStatus: %w", err)
	}

	// Create updated record with existing ID but new values
	updatedRecord := entities.NewMentalHealthRecordFromExisting(
		existingRecord.ID(),
		existingRecord.UserID(),
		happyLevel,
		energyLevel,
		command.Notes,
		status,
		existingRecord.CreatedAt(),
		existingRecord.UpdatedAt(),
		existingRecord.DeletedAt(),
	)

	// Update record
	if err := uc.recordRepo.Update(ctx, updatedRecord); err != nil {
		return nil, fmt.Errorf("uc.recordRepo.Update: %w", err)
	}

	return updatedRecord, nil
}

func (uc *MentalHealthRecordUseCaseImpl) Delete(ctx context.Context, command commands.DeleteMentalHealthRecordCommand) error {
	// Parse record ID
	recordID, err := value_objects.NewMentalHealthRecordIDFromString(command.ID)
	if err != nil {
		return fmt.Errorf("value_objects.NewMentalHealthRecordIDFromString: %w", err)
	}

	// Find existing record to check ownership
	existingRecord, err := uc.recordRepo.GetByID(ctx, recordID)
	if err != nil {
		return fmt.Errorf("uc.recordRepo.GetByID: %w", err)
	}

	// Check if user owns this record
	if existingRecord.UserID().String() != command.UserID {
		return fmt.Errorf("unauthorized: user does not own this record")
	}

	// Delete record
	if err := uc.recordRepo.Delete(ctx, recordID); err != nil {
		return fmt.Errorf("uc.recordRepo.Delete: %w", err)
	}

	return nil
}

func (uc *MentalHealthRecordUseCaseImpl) GetByID(ctx context.Context, id string, userID string) (*entities.MentalHealthRecord, error) {
	// Parse record ID
	recordID, err := value_objects.NewMentalHealthRecordIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewMentalHealthRecordIDFromString: %w", err)
	}

	// Find record
	record, err := uc.recordRepo.GetByID(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("uc.recordRepo.GetByID: %w", err)
	}

	// Check if user owns this record
	if record.UserID().String() != userID {
		return nil, fmt.Errorf("unauthorized: user does not own this record")
	}

	return record, nil
}

func (uc *MentalHealthRecordUseCaseImpl) GetByCondition(ctx context.Context, userID string, startedAt *string, endedAt *string) ([]*entities.MentalHealthRecord, error) {
	// Create search condition
	userIDVO, err := value_objects.NewUserIDFromString(userID)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewUserIDFromString: %w", err)
	}

	filter := &repositories.MentalHealthRecordFilter{
		UserID: userIDVO,
	}

	// Parse date range if provided - now supports ISO 8601 format with timezone
	if startedAt != nil {
		startTime, err := timeutil.ParseTime(*startedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format, expected ISO 8601 (e.g., 2025-08-23T17:00:00.000Z): %w", err)
		}
		filter.StartedAt = &startTime
	}

	if endedAt != nil {
		endTime, err := timeutil.ParseTime(*endedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format, expected ISO 8601 (e.g., 2025-08-23T17:00:00.000Z): %w", err)
		}
		filter.EndedAt = &endTime
	}

	// Get records from repository
	records, err := uc.recordRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("uc.recordRepo.GetByFilter: %w", err)
	}

	return records, nil
}

func (uc *MentalHealthRecordUseCaseImpl) GetHeatmap(ctx context.Context, command *commands.GetMentalHealthHeatmapCommand) (*commands.MentalHealthHeatmapResult, error) { // Create search condition
	userIDVO, err := value_objects.NewUserIDFromString(command.UserID)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewUserIDFromString: %w", err)
	}

	filter := &repositories.MentalHealthRecordFilter{
		UserID: userIDVO,
	}

	// Parse date range if provided - now supports ISO 8601 format with timezone
	if command.StartedAt != nil {
		startTime, err := timeutil.ParseTime(*command.StartedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format, expected ISO 8601 (e.g., 2025-08-23T17:00:00.000Z): %w", err)
		}
		filter.StartedAt = &startTime
	}

	if command.EndedAt != nil {
		endTime, err := timeutil.ParseTime(*command.EndedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format, expected ISO 8601 (e.g., 2025-08-23T17:00:00.000Z): %w", err)
		}
		filter.EndedAt = &endTime
	}

	// Get records from repository
	records, err := uc.recordRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("uc.recordRepo.GetByFilter: %w", err)
	}

	// Build heatmap data with proper calculation logic
	dateData := make(map[string]commands.HeatmapDataPoint)

	for _, record := range records {
		// Format date in UTC for consistent grouping
		date := record.CreatedAt().Format("2006-01-02")

		if existingData, exists := dateData[date]; exists {
			// Calculate running average for multiple records on same date
			totalHappy := existingData.HappyLevel*existingData.Count + record.HappyLevel().Value()
			totalEnergy := existingData.EnergyLevel*existingData.Count + record.EnergyLevel().Value()
			newCount := existingData.Count + 1

			dateData[date] = commands.HeatmapDataPoint{
				HappyLevel:  totalHappy / newCount,
				EnergyLevel: totalEnergy / newCount,
				Count:       newCount,
			}
		} else {
			// First record for this date
			dateData[date] = commands.HeatmapDataPoint{
				HappyLevel:  record.HappyLevel().Value(),
				EnergyLevel: record.EnergyLevel().Value(),
				Count:       1,
			}
		}
	}

	// Create structured result
	result := &commands.MentalHealthHeatmapResult{
		Data:         dateData,
		TotalRecords: len(records),
		DateRange: commands.DateRange{
			StartedAt: command.StartedAt,
			EndedAt:   command.EndedAt,
		},
	}

	return result, nil
}

func (uc *MentalHealthRecordUseCaseImpl) GetStreak(ctx context.Context, command *commands.GetMentalHealthStreakCommand) (*commands.MentalHealthStreakResult, error) {
	// Create user ID value object
	userIDVO, err := value_objects.NewUserIDFromString(command.UserID)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewUserIDFromString: %w", err)
	}

	// Get distinct dates for the user
	dates, err := uc.recordRepo.GetDistinctDatesForUser(ctx, userIDVO)
	if err != nil {
		return nil, fmt.Errorf("uc.recordRepo.GetDistinctDatesForUser: %w", err)
	}

	if len(dates) == 0 {
		return &commands.MentalHealthStreakResult{
			Streak:        0,
			LastEntryDate: nil,
		}, nil
	}

	// Calculate streak using the dates
	streak := calculateStreakFromDates(dates)

	// Return result with last entry date
	return &commands.MentalHealthStreakResult{
		Streak:        streak,
		LastEntryDate: &dates[0], // First element is most recent due to DESC order
	}, nil
}

// Helper function to calculate streak from dates
func calculateStreakFromDates(dates []string) int {
	if len(dates) == 0 {
		return 0
	}

	// Get current date in UTC using timeutil
	now := time.Now().UTC()
	today := timeutil.FormatDate(now)
	yesterday := timeutil.FormatDate(now.AddDate(0, 0, -1))

	// Parse the most recent entry date
	mostRecentDate := dates[0] // dates are in DESC order

	// Check if the most recent entry is today or yesterday
	// If not, streak is broken
	if mostRecentDate != today && mostRecentDate != yesterday {
		return 0
	}

	streak := 1

	for i := 1; i < len(dates); i++ {
		// Parse current date (more recent) using timeutil
		current, err := timeutil.ParseDate(dates[i-1])
		if err != nil {
			break
		}

		// Parse next date (going backwards in time) using timeutil
		next, err := timeutil.ParseDate(dates[i])
		if err != nil {
			break
		}

		// Calculate expected previous day
		expectedPrevDay := current.AddDate(0, 0, -1)

		// Check if dates are consecutive (next date should be 1 day before current)
		if timeutil.FormatDate(next) == timeutil.FormatDate(expectedPrevDay) {
			streak++
		} else {
			// Gap found, stop counting
			break
		}
	}

	return streak
}
