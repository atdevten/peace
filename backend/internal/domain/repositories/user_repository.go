package repositories

import (
	"context"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

type UserFilter struct {
	ID       *value_objects.UserID
	Email    *value_objects.Email
	Username *value_objects.Username
}

func NewUserFilter(id *value_objects.UserID, email *value_objects.Email, username *value_objects.Username) *UserFilter {
	return &UserFilter{
		ID:       id,
		Email:    email,
		Username: username,
	}
}

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id *value_objects.UserID) (*entities.User, error)
	GetByFilter(ctx context.Context, filter *UserFilter) (*entities.User, error)
	GetAll(ctx context.Context) ([]*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id *value_objects.UserID) error
	EmailExists(ctx context.Context, email value_objects.Email) (bool, error)
	UsernameExists(ctx context.Context, username value_objects.Username) (bool, error)
}
