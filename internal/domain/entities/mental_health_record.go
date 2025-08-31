package entities

import (
	"time"

	"github.com/atdevten/peace/internal/domain/value_objects"
)

type MentalHealthRecord struct {
	id          *value_objects.MentalHealthRecordID
	userID      *value_objects.UserID
	happyLevel  *value_objects.HappyLevel
	energyLevel *value_objects.EnergyLevel
	notes       *string
	status      *value_objects.MentalHealthRecordStatus
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
}

func NewMentalHealthRecord(
	userID string,
	happyLevel int,
	energyLevel int,
	notes *string,
	status string,
) (*MentalHealthRecord, error) {
	userIDVO, err := value_objects.NewUserIDFromString(userID)
	if err != nil {
		return nil, err
	}

	happyLevelVO, err := value_objects.NewHappyLevel(happyLevel)
	if err != nil {
		return nil, err
	}

	energyLevelVO, err := value_objects.NewEnergyLevel(energyLevel)
	if err != nil {
		return nil, err
	}

	statusVO, err := value_objects.NewMentalHealthRecordStatus(status)
	if err != nil {
		return nil, err
	}

	return &MentalHealthRecord{
		id:          value_objects.NewMentalHealthRecordID(),
		userID:      userIDVO,
		happyLevel:  happyLevelVO,
		energyLevel: energyLevelVO,
		notes:       notes,
		status:      statusVO,
	}, nil
}

func NewMentalHealthRecordFromExisting(
	id *value_objects.MentalHealthRecordID,
	userID *value_objects.UserID,
	happyLevel *value_objects.HappyLevel,
	energyLevel *value_objects.EnergyLevel,
	notes *string,
	status *value_objects.MentalHealthRecordStatus,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *MentalHealthRecord {
	return &MentalHealthRecord{
		id:          id,
		userID:      userID,
		happyLevel:  happyLevel,
		energyLevel: energyLevel,
		notes:       notes,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		deletedAt:   deletedAt,
	}
}

func (m *MentalHealthRecord) ID() *value_objects.MentalHealthRecordID {
	return m.id
}

func (m *MentalHealthRecord) UserID() *value_objects.UserID {
	return m.userID
}

func (m *MentalHealthRecord) HappyLevel() *value_objects.HappyLevel {
	return m.happyLevel
}

func (m *MentalHealthRecord) EnergyLevel() *value_objects.EnergyLevel {
	return m.energyLevel
}

func (m *MentalHealthRecord) Notes() *string {
	return m.notes
}

func (m *MentalHealthRecord) CreatedAt() time.Time {
	return m.createdAt
}

func (m *MentalHealthRecord) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *MentalHealthRecord) DeletedAt() *time.Time {
	return m.deletedAt
}

func (m *MentalHealthRecord) Status() *value_objects.MentalHealthRecordStatus {
	return m.status
}
