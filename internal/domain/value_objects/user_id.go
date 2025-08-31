package value_objects

import (
	"github.com/google/uuid"
)

type UserID struct {
	value string
}

func NewUserID() *UserID {
	return &UserID{value: uuid.New().String()}
}

func NewUserIDFromString(id string) (*UserID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, err
	}
	return &UserID{value: id}, nil
}

func (u *UserID) String() string {
	return u.value
}

func (u *UserID) IsZero() bool {
	return u.value == ""
}
