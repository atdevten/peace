package usecases

import (
	"context"
	"errors"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

type UserUseCase interface {
	GetByID(ctx context.Context, userID string) (*entities.User, error)
	UpdateProfile(ctx context.Context, userID string, firstName *string, lastName *string) (*entities.User, error)
	UpdatePassword(ctx context.Context, userID string, newPassword string) error
	Deactivate(ctx context.Context, userID string) error
	Delete(ctx context.Context, userID string) error
}

type UserUseCaseImpl struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) UserUseCase {
	return &UserUseCaseImpl{userRepo: userRepo}
}

func (uc *UserUseCaseImpl) GetByID(ctx context.Context, userID string) (*entities.User, error) {
	id, err := value_objects.NewUserIDFromString(userID)
	if err != nil {
		return nil, err
	}
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCaseImpl) UpdateProfile(ctx context.Context, userID string, firstName *string, lastName *string) (*entities.User, error) {
	user, err := uc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := user.UpdateProfile(firstName, lastName); err != nil {
		return nil, err
	}
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCaseImpl) UpdatePassword(ctx context.Context, userID string, newPassword string) error {
	if newPassword == "" {
		return errors.New("new password is required")
	}
	user, err := uc.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := user.UpdatePassword(newPassword); err != nil {
		return err
	}
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCaseImpl) Deactivate(ctx context.Context, userID string) error {
	user, err := uc.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := user.Deactivate(); err != nil {
		return err
	}
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCaseImpl) Delete(ctx context.Context, userID string) error {
	id, err := value_objects.NewUserIDFromString(userID)
	if err != nil {
		return err
	}
	return uc.userRepo.Delete(ctx, id)
}
