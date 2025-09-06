package entities

import (
	"time"

	"github.com/atdevten/peace/internal/domain/value_objects"
)

// UserOnlineStatus represents a user's online status
type UserOnlineStatus struct {
	id        *value_objects.UserOnlineStatusID
	userID    *value_objects.UserID
	userEmail *value_objects.Email
	isOnline  bool
	lastSeen  time.Time
	createdAt time.Time
	updatedAt time.Time
}

// NewUserOnlineStatus creates a new UserOnlineStatus entity
func NewUserOnlineStatus(userID *value_objects.UserID, userEmail *value_objects.Email) *UserOnlineStatus {
	now := time.Now()
	return &UserOnlineStatus{
		id:        value_objects.NewUserOnlineStatusID(),
		userID:    userID,
		userEmail: userEmail,
		isOnline:  true,
		lastSeen:  now,
		createdAt: now,
		updatedAt: now,
	}
}

// Getters
func (u *UserOnlineStatus) ID() *value_objects.UserOnlineStatusID {
	return u.id
}

func (u *UserOnlineStatus) UserID() *value_objects.UserID {
	return u.userID
}

func (u *UserOnlineStatus) UserEmail() *value_objects.Email {
	return u.userEmail
}

func (u *UserOnlineStatus) IsOnline() bool {
	return u.isOnline
}

func (u *UserOnlineStatus) LastSeen() time.Time {
	return u.lastSeen
}

func (u *UserOnlineStatus) CreatedAt() time.Time {
	return u.createdAt
}

func (u *UserOnlineStatus) UpdatedAt() time.Time {
	return u.updatedAt
}

// Business methods
func (u *UserOnlineStatus) GoOnline() {
	u.isOnline = true
	u.lastSeen = time.Now()
	u.updatedAt = time.Now()
}

func (u *UserOnlineStatus) GoOffline() {
	u.isOnline = false
	u.lastSeen = time.Now()
	u.updatedAt = time.Now()
}

func (u *UserOnlineStatus) UpdateLastSeen() {
	u.lastSeen = time.Now()
	u.updatedAt = time.Now()
}
