package value_objects

import (
	"fmt"
	"strings"
)

type MentalHealthRecordStatus string

const (
	RecordStatusPublic  MentalHealthRecordStatus = "public"
	RecordStatusPrivate MentalHealthRecordStatus = "private"
)

func (s MentalHealthRecordStatus) String() string {
	return string(s)
}

func NewMentalHealthRecordStatus(status string) (*MentalHealthRecordStatus, error) {
	status = strings.TrimSpace(status)

	if status != string(RecordStatusPublic) && status != string(RecordStatusPrivate) {
		return nil, fmt.Errorf("invalid mental health record status: %s", status)
	}

	statusVO := MentalHealthRecordStatus(status)
	return &statusVO, nil
}
