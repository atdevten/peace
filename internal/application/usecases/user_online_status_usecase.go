package usecases

import (
	"context"
	"fmt"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

// UserOnlineStatusUseCase defines the interface for user online status business logic
type UserOnlineStatusUseCase interface {
	// SetUserOnline sets a user as online
	SetUserOnline(ctx context.Context, userID, userEmail string) error

	// SetUserOffline sets a user as offline
	SetUserOffline(ctx context.Context, userID string) error

	// UpdateUserLastSeen updates the last seen timestamp for a user
	UpdateUserLastSeen(ctx context.Context, userID string) error

	// GetUserOnlineStatus gets the online status of a user
	GetUserOnlineStatus(ctx context.Context, userID string) (*entities.UserOnlineStatus, error)

	// GetOnlineUsers gets all online users
	GetOnlineUsers(ctx context.Context) ([]*entities.UserOnlineStatus, error)

	// IsUserOnline checks if a user is currently online
	IsUserOnline(ctx context.Context, userID string) (bool, error)

	// GetOnlineCount returns the number of online users quickly
	GetOnlineCount(ctx context.Context) (int64, error)
}

// UserOnlineStatusUseCaseImpl implements UserOnlineStatusUseCase
type UserOnlineStatusUseCaseImpl struct {
	userOnlineStatusRepo repositories.UserOnlineStatusRepository
}

// NewUserOnlineStatusUseCase creates a new user online status use case
func NewUserOnlineStatusUseCase(userOnlineStatusRepo repositories.UserOnlineStatusRepository) UserOnlineStatusUseCase {
	return &UserOnlineStatusUseCaseImpl{
		userOnlineStatusRepo: userOnlineStatusRepo,
	}
}

// SetUserOnline sets a user as online
func (uc *UserOnlineStatusUseCaseImpl) SetUserOnline(ctx context.Context, userID, userEmail string) error {
	// Create value objects
	userIDVO, err := value_objects.NewUserIDFromString(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	userEmailVO, err := value_objects.NewEmail(userEmail)
	if err != nil {
		return fmt.Errorf("invalid user email: %w", err)
	}

	// Create user online status entity
	status := entities.NewUserOnlineStatus(userIDVO, userEmailVO)
	status.GoOnline()

	// Save to repository
	err = uc.userOnlineStatusRepo.Save(ctx, status)
	if err != nil {
		return fmt.Errorf("failed to save user online status: %w", err)
	}

	return nil
}

// SetUserOffline sets a user as offline
func (uc *UserOnlineStatusUseCaseImpl) SetUserOffline(ctx context.Context, userID string) error {
	// Get current status
	status, err := uc.userOnlineStatusRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user status: %w", err)
	}

	if status != nil {
		// Update status to offline
		status.GoOffline()

		// Save updated status
		err = uc.userOnlineStatusRepo.Save(ctx, status)
		if err != nil {
			return fmt.Errorf("failed to save user offline status: %w", err)
		}
	}

	return nil
}

// UpdateUserLastSeen updates the last seen timestamp for a user
func (uc *UserOnlineStatusUseCaseImpl) UpdateUserLastSeen(ctx context.Context, userID string) error {
	err := uc.userOnlineStatusRepo.UpdateLastSeen(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to update user last seen: %w", err)
	}

	return nil
}

// GetUserOnlineStatus gets the online status of a user
func (uc *UserOnlineStatusUseCaseImpl) GetUserOnlineStatus(ctx context.Context, userID string) (*entities.UserOnlineStatus, error) {
	status, err := uc.userOnlineStatusRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user online status: %w", err)
	}

	return status, nil
}

// GetOnlineUsers gets all online users
func (uc *UserOnlineStatusUseCaseImpl) GetOnlineUsers(ctx context.Context) ([]*entities.UserOnlineStatus, error) {
	onlineUsers, err := uc.userOnlineStatusRepo.GetOnlineUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get online users: %w", err)
	}

	return onlineUsers, nil
}

// IsUserOnline checks if a user is currently online
func (uc *UserOnlineStatusUseCaseImpl) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	status, err := uc.GetUserOnlineStatus(ctx, userID)
	if err != nil {
		return false, err
	}

	if status == nil {
		return false, nil
	}

	return status.IsOnline(), nil
}

// GetOnlineCount returns the number of online users quickly
func (uc *UserOnlineStatusUseCaseImpl) GetOnlineCount(ctx context.Context) (int64, error) {
	n, err := uc.userOnlineStatusRepo.GetOnlineCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get online count: %w", err)
	}
	return n, nil
}
