package entities

import (
	"time"

	"github.com/atdevten/peace/internal/domain/value_objects"
)

type Quote struct {
	id        *value_objects.QuoteID
	content   *value_objects.Content
	author    *value_objects.Author
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewQuote(
	content string,
	author string,
) (*Quote, error) {
	contentVO, err := value_objects.NewContent(content)
	if err != nil {
		return nil, err
	}

	authorVO, err := value_objects.NewAuthor(author)
	if err != nil {
		return nil, err
	}

	return &Quote{
		id:      value_objects.NewQuoteID(),
		content: contentVO,
		author:  authorVO,
	}, nil
}

func NewQuoteFromExisting(
	id *value_objects.QuoteID,
	content *value_objects.Content,
	author *value_objects.Author,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *Quote {
	return &Quote{
		id:        id,
		content:   content,
		author:    author,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

func (q *Quote) ID() *value_objects.QuoteID {
	return q.id
}

func (q *Quote) Content() *value_objects.Content {
	return q.content
}

func (q *Quote) Author() *value_objects.Author {
	return q.author
}

func (q *Quote) CreatedAt() time.Time {
	return q.createdAt
}

func (q *Quote) UpdatedAt() time.Time {
	return q.updatedAt
}

func (q *Quote) DeletedAt() *time.Time {
	return q.deletedAt
}
