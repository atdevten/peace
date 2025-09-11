package models

import (
	"time"
)

type Quote struct {
	ID        int        `db:"id"`
	Content   string     `db:"content"`
	Author    string     `db:"author"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
