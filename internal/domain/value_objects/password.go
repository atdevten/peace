package value_objects

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	value string
}

type HashedPassword struct {
	hash string
}

func NewPassword(password string) (*Password, error) {
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	if len(password) > 72 { // bcrypt limitation
		return nil, errors.New("password too long")
	}

	// Check for at least one uppercase, one lowercase, one digit
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit {
		return nil, errors.New("password must contain at least one uppercase letter, one lowercase letter, and one digit")
	}

	return &Password{value: password}, nil
}

func (p *Password) Hash() (*HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &HashedPassword{hash: string(hash)}, nil
}

func NewHashedPassword(hash string) *HashedPassword {
	return &HashedPassword{hash: hash}
}

func (hp *HashedPassword) Verify(password *Password) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hp.hash), []byte(password.value))
	return err == nil
}

func (hp *HashedPassword) String() string {
	return hp.hash
}
