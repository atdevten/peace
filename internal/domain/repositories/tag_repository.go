package repositories

import (
	"context"
	"errors"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

var (
	ErrTagNotFound = errors.New("tag not found")
)

type TagFilter struct {
	ID   *value_objects.TagID
	Name *string
}

func NewTagFilter(id *value_objects.TagID, name *string) *TagFilter {
	return &TagFilter{
		ID:   id,
		Name: name,
	}
}

type TagRepository interface {
	Create(ctx context.Context, tag *entities.Tag) error
	GetByID(ctx context.Context, id *value_objects.TagID) (*entities.Tag, error)
	GetByFilter(ctx context.Context, filter *TagFilter) (*entities.Tag, error)
	GetAll(ctx context.Context) ([]*entities.Tag, error)
	Update(ctx context.Context, tag *entities.Tag) error
	Delete(ctx context.Context, id *value_objects.TagID) error
	GetByQuoteID(ctx context.Context, quoteID *value_objects.QuoteID) ([]*entities.Tag, error)
	AddTagToQuote(ctx context.Context, quoteID *value_objects.QuoteID, tagID *value_objects.TagID) error
	RemoveTagFromQuote(ctx context.Context, quoteID *value_objects.QuoteID, tagID *value_objects.TagID) error
}
