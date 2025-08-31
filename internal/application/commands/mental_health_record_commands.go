package commands

import (
	"errors"
)

type CreateMentalHealthRecordCommand struct {
	UserID      string
	HappyLevel  int
	EnergyLevel int
	Notes       *string
	Status      string
}

func NewCreateMentalHealthRecordCommand(userID string, happyLevel int, energyLevel int, notes *string, status string) (CreateMentalHealthRecordCommand, error) {
	if userID == "" {
		return CreateMentalHealthRecordCommand{}, errors.New("user_id is required")
	}

	return CreateMentalHealthRecordCommand{
		UserID:      userID,
		HappyLevel:  happyLevel,
		EnergyLevel: energyLevel,
		Notes:       notes,
		Status:      status,
	}, nil
}

type UpdateMentalHealthRecordCommand struct {
	ID          string
	UserID      string
	HappyLevel  int
	EnergyLevel int
	Notes       *string
	Status      string
}

func NewUpdateMentalHealthRecordCommand(id string, userID string, happyLevel int, energyLevel int, notes *string, status string) (UpdateMentalHealthRecordCommand, error) {
	if id == "" {
		return UpdateMentalHealthRecordCommand{}, errors.New("id is required")
	}

	if userID == "" {
		return UpdateMentalHealthRecordCommand{}, errors.New("user_id is required")
	}

	return UpdateMentalHealthRecordCommand{
		ID:          id,
		UserID:      userID,
		HappyLevel:  happyLevel,
		EnergyLevel: energyLevel,
		Notes:       notes,
		Status:      status,
	}, nil
}

type DeleteMentalHealthRecordCommand struct {
	ID     string
	UserID string
}

func NewDeleteMentalHealthRecordCommand(id string, userID string) (DeleteMentalHealthRecordCommand, error) {
	if id == "" {
		return DeleteMentalHealthRecordCommand{}, errors.New("id is required")
	}

	if userID == "" {
		return DeleteMentalHealthRecordCommand{}, errors.New("user_id is required")
	}

	return DeleteMentalHealthRecordCommand{
		ID:     id,
		UserID: userID,
	}, nil
}

type GetMentalHealthHeatmapCommand struct {
	UserID    string
	StartedAt *string
	EndedAt   *string
}

func NewGetMentalHealthHeatmapCommand(userID string, startedAt *string, endedAt *string) *GetMentalHealthHeatmapCommand {
	return &GetMentalHealthHeatmapCommand{
		UserID:    userID,
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}
}

// Application layer response structs
type HeatmapDataPoint struct {
	HappyLevel  int
	EnergyLevel int
	Count       int
}

type DateRange struct {
	StartedAt *string
	EndedAt   *string
}

type MentalHealthHeatmapResult struct {
	Data         map[string]HeatmapDataPoint
	TotalRecords int
	DateRange    DateRange
}

type GetMentalHealthStreakCommand struct {
	UserID string
}

func NewGetMentalHealthStreakCommand(userID string) *GetMentalHealthStreakCommand {
	return &GetMentalHealthStreakCommand{
		UserID: userID,
	}
}

type MentalHealthStreakResult struct {
	Streak        int     `json:"streak"`
	LastEntryDate *string `json:"last_entry_date,omitempty"`
}
