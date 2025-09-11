package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
	"github.com/atdevten/peace/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type PostgreSQLUserRepository struct {
	db *gorm.DB
}

func NewPostgreSQLUserRepository(db *gorm.DB) repositories.UserRepository {
	return &PostgreSQLUserRepository{
		db: db,
	}
}

func (r *PostgreSQLUserRepository) Create(ctx context.Context, user *entities.User) error {
	// Convert value objects to strings for database storage
	var firstName *string
	if user.FirstName() != nil {
		firstNameStr := user.FirstName().String()
		firstName = &firstNameStr
	}

	var lastName *string
	if user.LastName() != nil {
		lastNameStr := user.LastName().String()
		lastName = &lastNameStr
	}

	var passwordHash string
	if user.PasswordHash() != nil {
		passwordHash = user.PasswordHash().String()
	}

	model := models.User{
		ID:            user.ID().String(),
		Email:         user.Email().String(),
		Username:      user.Username().String(),
		FirstName:     firstName,
		LastName:      lastName,
		PasswordHash:  passwordHash,
		IsActive:      user.IsActive(),
		EmailVerified: user.EmailVerified(),
		AuthProvider:  user.AuthProvider(),
		GoogleID:      user.GoogleID(),
		GooglePicture: user.GooglePicture(),
		CreatedAt:     user.CreatedAt(),
		UpdatedAt:     user.UpdatedAt(),
		DeletedAt:     user.DeletedAt(),
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("r.db.Create: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) GetByID(ctx context.Context, id *value_objects.UserID) (*entities.User, error) {
	var model models.User

	result := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("r.db.First: %w", result.Error)
	}

	return r.modelToEntity(model)
}

func (r *PostgreSQLUserRepository) GetByFilter(ctx context.Context, filter *repositories.UserFilter) (*entities.User, error) {
	var model models.User
	var query string
	var queryValue string

	if filter.ID != nil {
		query = "id = ?"
		queryValue = filter.ID.String()
	} else if filter.Email != nil {
		query = "email = ?"
		queryValue = filter.Email.String()
	} else if filter.Username != nil {
		query = "username = ?"
		queryValue = filter.Username.String()
	} else {
		return nil, fmt.Errorf("no search criteria provided")
	}

	result := r.db.WithContext(ctx).Where(query, queryValue).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("r.db.First: %w", result.Error)
	}

	return r.modelToEntity(model)
}

func (r *PostgreSQLUserRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	var models []models.User

	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("r.db.Find: %w", err)
	}

	var users []*entities.User
	for _, model := range models {
		user, err := r.modelToEntity(model)
		if err != nil {
			return nil, fmt.Errorf("modelToEntity: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *PostgreSQLUserRepository) Update(ctx context.Context, user *entities.User) error {
	// Convert value objects to strings for database storage
	var firstName *string
	if user.FirstName() != nil {
		firstNameStr := user.FirstName().String()
		firstName = &firstNameStr
	}

	var lastName *string
	if user.LastName() != nil {
		lastNameStr := user.LastName().String()
		lastName = &lastNameStr
	}

	var passwordHash string
	if user.PasswordHash() != nil {
		passwordHash = user.PasswordHash().String()
	}

	model := models.User{
		ID:            user.ID().String(),
		Email:         user.Email().String(),
		Username:      user.Username().String(),
		FirstName:     firstName,
		LastName:      lastName,
		PasswordHash:  passwordHash,
		IsActive:      user.IsActive(),
		EmailVerified: user.EmailVerified(),
		AuthProvider:  user.AuthProvider(),
		GoogleID:      user.GoogleID(),
		GooglePicture: user.GooglePicture(),
		CreatedAt:     user.CreatedAt(),
		UpdatedAt:     user.UpdatedAt(),
		DeletedAt:     user.DeletedAt(),
	}

	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", model.ID).Updates(&model).Error; err != nil {
		return fmt.Errorf("r.db.Updates: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) Delete(ctx context.Context, id *value_objects.UserID) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id.String()).Error; err != nil {
		return fmt.Errorf("r.db.Delete: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) EmailExists(ctx context.Context, email value_objects.Email) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email.String()).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("r.db.Count: %w", err)
	}
	return count > 0, nil
}

func (r *PostgreSQLUserRepository) UsernameExists(ctx context.Context, username value_objects.Username) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username.String()).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("r.db.Count: %w", err)
	}
	return count > 0, nil
}

// Helper method to convert model to entity
func (r *PostgreSQLUserRepository) modelToEntity(model models.User) (*entities.User, error) {
	userID, err := value_objects.NewUserIDFromString(model.ID)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewUserIDFromString: %w", err)
	}

	email, err := value_objects.NewEmail(model.Email)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewEmail: %w", err)
	}

	username, err := value_objects.NewUsername(model.Username)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewUsername: %w", err)
	}

	var hashedPassword *value_objects.HashedPassword
	if model.PasswordHash != "" {
		hashedPassword = value_objects.NewHashedPassword(model.PasswordHash)
	}

	// Convert database strings to value objects
	firstName, err := value_objects.NewOptionalFirstName(model.FirstName)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewOptionalFirstName: %w", err)
	}

	lastName, err := value_objects.NewOptionalLastName(model.LastName)
	if err != nil {
		return nil, fmt.Errorf("value_objects.NewOptionalLastName: %w", err)
	}

	return entities.NewUserFromRepository(
		userID,
		email,
		username,
		firstName,
		lastName,
		hashedPassword,
		model.IsActive,
		model.EmailVerified,
		model.AuthProvider,
		model.GoogleID,
		model.GooglePicture,
		model.CreatedAt,
		model.UpdatedAt,
		model.DeletedAt,
	), nil
}
