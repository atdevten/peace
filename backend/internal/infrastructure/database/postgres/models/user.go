package models

import (
	"time"
)

type User struct {
	ID            string     `db:"id"`
	Email         string     `db:"email"`
	Username      string     `db:"username"`
	FirstName     *string    `db:"first_name"`
	LastName      *string    `db:"last_name"`
	PasswordHash  string     `db:"password_hash"`
	IsActive      bool       `db:"is_active"`
	EmailVerified bool       `db:"email_verified"`
	AuthProvider  string     `db:"auth_provider"`
	GoogleID      *string    `db:"google_id"`
	GooglePicture *string    `db:"google_picture"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}
