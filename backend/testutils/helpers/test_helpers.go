package helpers

import (
	"time"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int
func IntPtr(i int) *int {
	return &i
}

// TimePtr returns a pointer to the given time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// CreateTestUser creates a test user with default values
func CreateTestUser() *entities.User {
	user, _ := entities.NewUser(
		"test@example.com",
		"testuser",
		StringPtr("John"),
		StringPtr("Doe"),
		"Password123",
	)
	return user
}

// CreateTestGoogleUser creates a test Google user
func CreateTestGoogleUser() *entities.User {
	user, _ := entities.NewGoogleUser(
		"google@example.com",
		StringPtr("Google"),
		StringPtr("User"),
		"google123",
		StringPtr("https://example.com/avatar.jpg"),
	)
	return user
}

// CreateTestUserID creates a test user ID
func CreateTestUserID() *value_objects.UserID {
	return value_objects.NewUserID()
}

// CreateTestEmail creates a test email
func CreateTestEmail(email string) *value_objects.Email {
	emailVO, _ := value_objects.NewEmail(email)
	return emailVO
}

// CreateTestUsername creates a test username
func CreateTestUsername(username string) *value_objects.Username {
	usernameVO, _ := value_objects.NewUsername(username)
	return usernameVO
}

// CreateTestPassword creates a test password
func CreateTestPassword(password string) *value_objects.Password {
	passwordVO, _ := value_objects.NewPassword(password)
	return passwordVO
}

// CreateTestMentalHealthRecord creates a test mental health record
func CreateTestMentalHealthRecord() *entities.MentalHealthRecord {
	record, _ := entities.NewMentalHealthRecord(
		"550e8400-e29b-41d4-a716-446655440000", // Fixed user ID for testing
		5,                                      // happy level
		7,                                      // energy level
		StringPtr("Feeling good today"),
		"public", // status
	)
	return record
}

// CreateTestMentalHealthRecordWithUserID creates a test mental health record with specific user ID
func CreateTestMentalHealthRecordWithUserID(userID *value_objects.UserID) *entities.MentalHealthRecord {
	record, _ := entities.NewMentalHealthRecord(
		userID.String(), // Use provided user ID
		5,               // happy level
		7,               // energy level
		StringPtr("Feeling good today"),
		"public", // status
	)
	return record
}

// CreateTestQuote creates a test quote
func CreateTestQuote() *entities.Quote {
	quote, _ := entities.NewQuote(
		"Life is beautiful",
		"Anonymous",
	)
	return quote
}

// CreateTestTag creates a test tag
func CreateTestTag() *entities.Tag {
	tag, _ := entities.NewTag("motivation", "Motivational quotes and content")
	return tag
}
