package models

import "time"

type MentalHealthRecord struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	HappyLevel  int        `db:"happy_level"`
	EnergyLevel int        `db:"energy_level"`
	Notes       *string    `db:"notes"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
