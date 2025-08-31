package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/atdevten/peace/internal/application/commands"
	appjwt "github.com/atdevten/peace/internal/application/services/jwt"
	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

type AuthUseCase interface {
	Register(ctx context.Context, command commands.RegisterCommand) (*entities.User, error)
	Login(ctx context.Context, command commands.LoginCommand) (*entities.User, string, string, error) // user, access, refresh, error
	Refresh(ctx context.Context, accessToken string, refreshToken string) (string, string, error)     // new access, new refresh, error
}

type AuthUseCaseImpl struct {
	userRepo   repositories.UserRepository
	jwtService appjwt.Service
}

func NewAuthUseCase(userRepo repositories.UserRepository, jwtService appjwt.Service) AuthUseCase {
	return &AuthUseCaseImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *AuthUseCaseImpl) Register(ctx context.Context, command commands.RegisterCommand) (*entities.User, error) {
	// Create new user
	user, err := entities.NewUser(
		command.Email,
		command.Username,
		command.FirstName,
		command.LastName,
		command.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("entities.NewUser: %w", err)
	}

	// Check if email already exists
	emailExists, err := uc.userRepo.EmailExists(ctx, *user.Email())
	if err != nil {
		return nil, fmt.Errorf("uc.userRepo.EmailExists: %w", err)
	}
	if emailExists {
		return nil, errors.New("email already registered")
	}

	// Create user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("uc.userRepo.Create: %w", err)
	}

	return user, nil
}

func (uc *AuthUseCaseImpl) Login(ctx context.Context, command commands.LoginCommand) (*entities.User, string, string, error) {
	// Find user by email
	emailVO, err := value_objects.NewEmail(command.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("value_objects.NewEmail: %w", err)
	}

	user, err := uc.userRepo.GetByFilter(ctx, repositories.NewUserFilter(nil, emailVO, nil))
	if err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	// Check if user can login
	if err = user.CanLogin(); err != nil {
		return nil, "", "", fmt.Errorf("user.CanLogin: %w", err)
	}

	// Verify password
	if err = user.VerifyPassword(command.Password); err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	// Generate tokens
	access, err := uc.jwtService.GenerateAccessToken(*user.ID(), *user.Email())
	if err != nil {
		return nil, "", "", fmt.Errorf("uc.jwtService.GenerateAccessToken: %w", err)
	}

	refresh, err := uc.jwtService.GenerateRefreshToken(*user.ID(), *user.Email())
	if err != nil {
		return nil, "", "", fmt.Errorf("uc.jwtService.GenerateRefreshToken: %w", err)
	}

	return user, access, refresh, nil
}

func (uc *AuthUseCaseImpl) Refresh(ctx context.Context, accessToken string, refreshToken string) (string, string, error) {
	// Validate refresh token first
	refreshClaims, err := uc.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Optionally: check access token expired/near expiry via ValidateAccessToken
	// For simplicity, issue new tokens based on refresh claims
	userID, err := value_objects.NewUserIDFromString(refreshClaims.UserID)
	if err != nil {
		return "", "", err
	}
	email, err := value_objects.NewEmail(refreshClaims.Email)
	if err != nil {
		return "", "", err
	}

	newAccess, err := uc.jwtService.GenerateAccessToken(*userID, *email)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := uc.jwtService.GenerateRefreshToken(*userID, *email)
	if err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}
