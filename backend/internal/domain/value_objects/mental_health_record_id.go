package value_objects

import "github.com/google/uuid"

type MentalHealthRecordID struct {
	value string
}

func NewMentalHealthRecordID() *MentalHealthRecordID {
	return &MentalHealthRecordID{value: uuid.New().String()}
}

func NewMentalHealthRecordIDFromString(id string) (*MentalHealthRecordID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, err
	}
	return &MentalHealthRecordID{value: id}, nil
}

func (m *MentalHealthRecordID) String() string {
	return m.value
}

func (m *MentalHealthRecordID) IsZero() bool {
	return m.value == ""
}
