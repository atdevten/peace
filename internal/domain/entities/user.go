package entities

import (
	"errors"
	"time"

	"github.com/atdevten/peace/internal/domain/value_objects"
)

// User represents a user in the domain
type User struct {
	id            *value_objects.UserID
	email         *value_objects.Email
	username      *value_objects.Username
	firstName     *value_objects.FirstName
	lastName      *value_objects.LastName
	passwordHash  *value_objects.HashedPassword
	isActive      bool
	emailVerified bool
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     *time.Time
}

// NewUser creates a new User entity with validation
func NewUser(
	email string,
	username string,
	firstName *string,
	lastName *string,
	password string,
) (*User, error) {
	// Create and validate value objects
	emailVO, err := value_objects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	usernameVO, err := value_objects.NewUsername(username)
	if err != nil {
		return nil, err
	}

	passwordVO, err := value_objects.NewPassword(password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := passwordVO.Hash()
	if err != nil {
		return nil, err
	}

	// Convert optional strings to value objects
	firstNameVO, err := value_objects.NewOptionalFirstName(firstName)
	if err != nil {
		return nil, err
	}

	lastNameVO, err := value_objects.NewOptionalLastName(lastName)
	if err != nil {
		return nil, err
	}

	return &User{
		id:            value_objects.NewUserID(),
		email:         emailVO,
		username:      usernameVO,
		firstName:     firstNameVO,
		lastName:      lastNameVO,
		passwordHash:  hashedPassword,
		isActive:      true,
		emailVerified: false,
	}, nil
}

// Getters
func (u *User) ID() *value_objects.UserID {
	return u.id
}

func (u *User) Email() *value_objects.Email {
	return u.email
}

func (u *User) Username() *value_objects.Username {
	return u.username
}

func (u *User) FirstName() *value_objects.FirstName {
	return u.firstName
}

func (u *User) LastName() *value_objects.LastName {
	return u.lastName
}

func (u *User) PasswordHash() *value_objects.HashedPassword {
	return u.passwordHash
}

func (u *User) IsActive() bool {
	return u.isActive
}

func (u *User) EmailVerified() bool {
	return u.emailVerified
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) DeletedAt() *time.Time {
	return u.deletedAt
}

// Business methods
func (u *User) VerifyPassword(password string) error {
	if !u.isActive {
		return errors.New("user account is deactivated")
	}

	passwordVO, err := value_objects.NewPassword(password)
	if err != nil {
		return err
	}

	if !u.passwordHash.Verify(passwordVO) {
		return errors.New("invalid password")
	}

	return nil
}

func (u *User) VerifyEmail() {
	u.emailVerified = true
	u.updatedAt = time.Now()
}

func (u *User) Deactivate() error {
	if !u.isActive {
		return errors.New("user is already deactivated")
	}

	u.isActive = false
	u.updatedAt = time.Now()
	return nil
}

func (u *User) Activate() error {
	if u.isActive {
		return errors.New("user is already active")
	}

	u.isActive = true
	u.updatedAt = time.Now()
	return nil
}

func (u *User) UpdatePassword(newPassword string) error {
	passwordVO, err := value_objects.NewPassword(newPassword)
	if err != nil {
		return err
	}

	hashedPassword, err := passwordVO.Hash()
	if err != nil {
		return err
	}

	u.passwordHash = hashedPassword
	u.updatedAt = time.Now()
	return nil
}

func (u *User) UpdateProfile(firstName *string, lastName *string) error {
	// Convert optional strings to value objects
	firstNameVO, err := value_objects.NewOptionalFirstName(firstName)
	if err != nil {
		return err
	}

	lastNameVO, err := value_objects.NewOptionalLastName(lastName)
	if err != nil {
		return err
	}

	u.firstName = firstNameVO
	u.lastName = lastNameVO
	u.updatedAt = time.Now()
	return nil
}

func (u *User) SoftDelete() error {
	if u.deletedAt != nil {
		return errors.New("user is already deleted")
	}

	now := time.Now()
	u.deletedAt = &now
	u.isActive = false
	u.updatedAt = now
	return nil
}

// Factory method from repository data
func NewUserFromRepository(
	id *value_objects.UserID,
	email *value_objects.Email,
	username *value_objects.Username,
	firstName *value_objects.FirstName,
	lastName *value_objects.LastName,
	passwordHash *value_objects.HashedPassword,
	isActive bool,
	emailVerified bool,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *User {
	return &User{
		id:            id,
		email:         email,
		username:      username,
		firstName:     firstName,
		lastName:      lastName,
		passwordHash:  passwordHash,
		isActive:      isActive,
		emailVerified: emailVerified,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		deletedAt:     deletedAt,
	}
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.firstName == nil && u.lastName == nil {
		return ""
	}

	firstName := ""
	if u.firstName != nil {
		firstName = u.firstName.Value()
	}

	lastName := ""
	if u.lastName != nil {
		lastName = u.lastName.Value()
	}

	return firstName + " " + lastName
}

// CanLogin checks if user can login
func (u *User) CanLogin() error {
	if !u.isActive {
		return errors.New("account is deactivated")
	}

	if u.deletedAt != nil {
		return errors.New("account is deleted")
	}

	return nil
}
