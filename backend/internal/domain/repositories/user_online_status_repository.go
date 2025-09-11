package repositories

import (
	"context"

	"github.com/atdevten/peace/internal/domain/entities"
)

// UserOnlineStatusRepository defines the interface for user online status persistence
type UserOnlineStatusRepository interface {
	// Save saves or updates a user online status
	Save(ctx context.Context, status *entities.UserOnlineStatus) error

	// GetByUserID retrieves online status by user ID
	GetByUserID(ctx context.Context, userID string) (*entities.UserOnlineStatus, error)

	// GetOnlineUsers retrieves all online users
	GetOnlineUsers(ctx context.Context) ([]*entities.UserOnlineStatus, error)

	// GetOnlineCount returns the number of online users quickly
	GetOnlineCount(ctx context.Context) (int64, error)

	// DeleteByUserID deletes online status by user ID
	DeleteByUserID(ctx context.Context, userID string) error

	// UpdateLastSeen updates the last seen timestamp for a user
	UpdateLastSeen(ctx context.Context, userID string) error
}
