package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
	redisclient "github.com/atdevten/peace/internal/infrastructure/database/redis"
)

type userOnlineDTO struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	IsOnline  bool   `json:"is_online"`
	LastSeen  int64  `json:"last_seen"`
}

// RedisUserOnlineStatusRepository implements UserOnlineStatusRepository using Redis
type RedisUserOnlineStatusRepository struct {
	client redisclient.Client
}

// NewRedisUserOnlineStatusRepository creates a new Redis repository
func NewRedisUserOnlineStatusRepository(client redisclient.Client) repositories.UserOnlineStatusRepository {
	return &RedisUserOnlineStatusRepository{
		client: client,
	}
}

// Save saves or updates a user online status
func (r *RedisUserOnlineStatusRepository) Save(ctx context.Context, status *entities.UserOnlineStatus) error {
	key := fmt.Sprintf("user:online:%s", status.UserID().String())

	dto := userOnlineDTO{
		UserID:    status.UserID().String(),
		UserEmail: status.UserEmail().String(),
		IsOnline:  status.IsOnline(),
		LastSeen:  time.Now().Unix(),
	}
	jsonData, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal user online dto: %w", err)
	}

	// Short TTL so that missing pings drop the user from online quickly
	expiration := 20 * time.Second
	if err := r.client.Set(ctx, key, string(jsonData), expiration); err != nil {
		return fmt.Errorf("failed to save user online dto to Redis: %w", err)
	}

	if status.IsOnline() {
		if err := r.client.SAdd(ctx, "users:online", status.UserID().String()); err != nil {
			return fmt.Errorf("failed to add user to online set: %w", err)
		}
	} else {
		if err := r.client.SRem(ctx, "users:online", status.UserID().String()); err != nil {
			return fmt.Errorf("failed to remove user from online set: %w", err)
		}
	}
	return nil
}

// GetByUserID retrieves online status by user ID
func (r *RedisUserOnlineStatusRepository) GetByUserID(ctx context.Context, userID string) (*entities.UserOnlineStatus, error) {
	key := fmt.Sprintf("user:online:%s", userID)

	exists, err := r.client.Exists(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to check key existence: %w", err)
	}
	if exists == 0 {
		return nil, nil
	}

	raw, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get user online dto: %w", err)
	}
	var dto userOnlineDTO
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dto: %w", err)
	}

	uid, err := value_objects.NewUserIDFromString(dto.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in dto: %w", err)
	}
	email, err := value_objects.NewEmail(dto.UserEmail)
	if err != nil {
		return nil, fmt.Errorf("invalid user email in dto: %w", err)
	}
	ent := entities.NewUserOnlineStatus(uid, email)
	if !dto.IsOnline {
		ent.GoOffline()
	}
	// lastSeen is managed internally; dto value is used for TTL/cleanup logic
	return ent, nil
}

// GetOnlineUsers retrieves all online users
func (r *RedisUserOnlineStatusRepository) GetOnlineUsers(ctx context.Context) ([]*entities.UserOnlineStatus, error) {
	ids, err := r.client.SMembers(ctx, "users:online")
	if err != nil {
		return nil, fmt.Errorf("failed to get online user IDs from Redis: %w", err)
	}
	if len(ids) == 0 {
		return []*entities.UserOnlineStatus{}, nil
	}

	// Build keys and bulk fetch via MGET for better performance
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf("user:online:%s", id))
	}

	rawVals, err := r.client.MGet(ctx, keys...)
	if err != nil {
		return nil, fmt.Errorf("failed to MGET online users: %w", err)
	}

	var results []*entities.UserOnlineStatus
	for i, raw := range rawVals {
		id := ids[i]
		if raw == "" {
			// key expired -> cleanup set membership
			_ = r.client.SRem(ctx, "users:online", id)
			continue
		}
		var dto userOnlineDTO
		if err := json.Unmarshal([]byte(raw), &dto); err != nil {
			continue
		}
		uid, err := value_objects.NewUserIDFromString(dto.UserID)
		if err != nil {
			continue
		}
		email, err := value_objects.NewEmail(dto.UserEmail)
		if err != nil {
			continue
		}
		ent := entities.NewUserOnlineStatus(uid, email)
		if !dto.IsOnline {
			ent.GoOffline()
		}
		results = append(results, ent)
	}
	return results, nil
}

// DeleteByUserID deletes online status by user ID
func (r *RedisUserOnlineStatusRepository) DeleteByUserID(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:online:%s", userID)
	_ = r.client.SRem(ctx, "users:online", userID)
	if err := r.client.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to delete user online status from Redis: %w", err)
	}
	return nil
}

// UpdateLastSeen updates the last seen timestamp for a user
func (r *RedisUserOnlineStatusRepository) UpdateLastSeen(ctx context.Context, userID string) error {
	st, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user status: %w", err)
	}
	if st == nil {
		return nil
	}
	st.UpdateLastSeen()
	if err := r.Save(ctx, st); err != nil {
		return fmt.Errorf("failed to save updated status: %w", err)
	}
	return nil
}

// GetOnlineCount returns the number of online users quickly via SCARD
func (r *RedisUserOnlineStatusRepository) GetOnlineCount(ctx context.Context) (int64, error) {
	count, err := r.client.SCard(ctx, "users:online")
	if err != nil {
		return 0, fmt.Errorf("failed to SCARD users:online: %w", err)
	}
	return count, nil
}
