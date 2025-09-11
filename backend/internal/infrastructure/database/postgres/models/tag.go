package models

import (
	"time"

	"github.com/atdevten/peace/internal/domain/value_objects"
	"gorm.io/gorm"
)

type Tag struct {
	ID          int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (t *Tag) TableName() string {
	return "tags"
}

// ToDomain converts GORM model to domain entity
func (t *Tag) ToDomain() (*value_objects.TagID, *value_objects.TagName, error) {
	tagID := value_objects.NewTagIDFromInt(t.ID)

	tagName, err := value_objects.NewTagName(t.Name)
	if err != nil {
		return nil, nil, err
	}

	return tagID, tagName, nil
}

// FromDomain converts domain entity to GORM model
func (t *Tag) FromDomain(tagID *value_objects.TagID, tagName *value_objects.TagName) {
	t.ID = tagID.IntValue()
	t.Name = tagName.Value()
}
