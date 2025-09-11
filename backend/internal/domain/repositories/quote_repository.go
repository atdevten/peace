package repositories

import (
	"context"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

type QuoteFilter struct {
	ID      *value_objects.QuoteID
	Author  *string
	Content *string
}

func NewQuoteFilter(id *value_objects.QuoteID, author *string, content *string) *QuoteFilter {
	return &QuoteFilter{
		ID:      id,
		Author:  author,
		Content: content,
	}
}

type QuoteRepository interface {
	Create(ctx context.Context, quote *entities.Quote) error
	GetByID(ctx context.Context, id *value_objects.QuoteID) (*entities.Quote, error)
	GetByFilter(ctx context.Context, filter *QuoteFilter) ([]*entities.Quote, error)
	GetAll(ctx context.Context) ([]*entities.Quote, error)
	GetRandom(ctx context.Context) (*entities.Quote, error)
	Update(ctx context.Context, quote *entities.Quote) error
	Delete(ctx context.Context, id *value_objects.QuoteID) error
}
