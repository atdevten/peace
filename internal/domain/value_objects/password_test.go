package value_objects

import (
	"testing"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "valid password",
			input:       "Password123",
			wantValue:   "Password123",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid password with special characters",
			input:       "MyP@ssw0rd",
			wantValue:   "MyP@ssw0rd",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "valid password minimum length",
			input:       "Pass1234",
			wantValue:   "Pass1234",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "password too short",
			input:       "Pass1",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "password must be at least 8 characters",
		},
		{
			name:        "password too long",
			input:       "VeryLongPasswordThatExceedsTheMaximumAllowedLengthForBcryptHashingAlgorithm",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "password too long",
		},
		{
			name:        "password missing uppercase",
			input:       "password123",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "password must contain at least one uppercase letter, one lowercase letter, and one digit",
		},
		{
			name:        "password missing lowercase",
			input:       "PASSWORD123",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "password must contain at least one uppercase letter, one lowercase letter, and one digit",
		},
		{
			name:        "password missing digit",
			input:       "Password",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "password must contain at least one uppercase letter, one lowercase letter, and one digit",
		},
		{
			name:        "password missing uppercase and digit",
			input:       "password",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "password must contain at least one uppercase letter, one lowercase letter, and one digit",
		},
		{
			name:        "password missing lowercase and digit",
			input:       "PASSWORD",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "password must contain at least one uppercase letter, one lowercase letter, and one digit",
		},
		{
			name:        "password missing uppercase and lowercase",
			input:       "12345678",
			wantValue:   "",
			wantErr:     true,
			expectedErr: "password must contain at least one uppercase letter, one lowercase letter, and one digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPassword(tt.input)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewPassword() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewPassword() error = %v, expected %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewPassword() unexpected error = %v", err)
				return
			}

			// Check value
			if got.value != tt.wantValue {
				t.Errorf("NewPassword() = %v, want %v", got.value, tt.wantValue)
			}
		})
	}
}

func TestPassword_Hash(t *testing.T) {
	password, err := NewPassword("Password123")
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	hashedPassword, err := password.Hash()
	if err != nil {
		t.Errorf("Password.Hash() unexpected error = %v", err)
		return
	}

	if hashedPassword.hash == "" {
		t.Error("Password.Hash() returned empty hash")
	}

	if hashedPassword.hash == password.value {
		t.Error("Password.Hash() returned unhashed password")
	}
}

func TestNewHashedPassword(t *testing.T) {
	hash := "$2a$10$abcdefghijklmnopqrstuvwxyz123456"
	hashedPassword := NewHashedPassword(hash)

	if hashedPassword.hash != hash {
		t.Errorf("NewHashedPassword() = %v, want %v", hashedPassword.hash, hash)
	}
}

func TestHashedPassword_Verify(t *testing.T) {
	// Create a password and hash it
	password, err := NewPassword("Password123")
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	hashedPassword, err := password.Hash()
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	if !hashedPassword.Verify(password) {
		t.Error("HashedPassword.Verify() failed for correct password")
	}

	// Test incorrect password
	wrongPassword, err := NewPassword("WrongPass123")
	if err != nil {
		t.Fatalf("Failed to create wrong password: %v", err)
	}

	if hashedPassword.Verify(wrongPassword) {
		t.Error("HashedPassword.Verify() succeeded for wrong password")
	}
}

func TestHashedPassword_String(t *testing.T) {
	hash := "$2a$10$abcdefghijklmnopqrstuvwxyz123456"
	hashedPassword := NewHashedPassword(hash)

	if hashedPassword.String() != hash {
		t.Errorf("HashedPassword.String() = %v, want %v", hashedPassword.String(), hash)
	}
}
