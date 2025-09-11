package models

import (
	"time"

	"github.com/atdevten/peace/internal/domain/value_objects"
)

type QuoteTag struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	QuoteID   int       `gorm:"not null;index" json:"quote_id"`
	TagID     int       `gorm:"not null;index" json:"tag_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (qt *QuoteTag) TableName() string {
	return "quote_tags"
}

// ToDomain converts GORM model to domain value objects
func (qt *QuoteTag) ToDomain() (*value_objects.QuoteID, *value_objects.TagID) {
	quoteID := value_objects.NewQuoteIDFromInt(qt.QuoteID)
	tagID := value_objects.NewTagIDFromInt(qt.TagID)
	return quoteID, tagID
}

// FromDomain converts domain value objects to GORM model
func (qt *QuoteTag) FromDomain(quoteID *value_objects.QuoteID, tagID *value_objects.TagID) {
	qt.QuoteID = quoteID.Value()
	qt.TagID = tagID.IntValue()
}
