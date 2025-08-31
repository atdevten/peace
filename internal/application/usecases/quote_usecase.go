package usecases

import (
	"context"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

type QuoteUseCase interface {
	CreateQuote(ctx context.Context, content, author string) error
	GetQuoteByID(ctx context.Context, id string) (*entities.Quote, error)
	GetAllQuotes(ctx context.Context) ([]*entities.Quote, error)
	GetRandomQuote(ctx context.Context) (*entities.Quote, error)
	GetQuotesByFilter(ctx context.Context, author *string, content *string) ([]*entities.Quote, error)
	UpdateQuote(ctx context.Context, id string, content, author string) error
	DeleteQuote(ctx context.Context, id string) error
}

type QuoteUseCaseImpl struct {
	quoteRepo repositories.QuoteRepository
}

func NewQuoteUseCase(quoteRepo repositories.QuoteRepository) QuoteUseCase {
	return &QuoteUseCaseImpl{
		quoteRepo: quoteRepo,
	}
}

func (u *QuoteUseCaseImpl) CreateQuote(ctx context.Context, content, author string) error {
	quote, err := entities.NewQuote(content, author)
	if err != nil {
		return err
	}

	return u.quoteRepo.Create(ctx, quote)
}

func (u *QuoteUseCaseImpl) GetQuoteByID(ctx context.Context, id string) (*entities.Quote, error) {
	quoteID, err := value_objects.NewQuoteIDFromString(id)
	if err != nil {
		return nil, err
	}

	return u.quoteRepo.GetByID(ctx, quoteID)
}

func (u *QuoteUseCaseImpl) GetAllQuotes(ctx context.Context) ([]*entities.Quote, error) {
	return u.quoteRepo.GetAll(ctx)
}

func (u *QuoteUseCaseImpl) GetRandomQuote(ctx context.Context) (*entities.Quote, error) {
	return u.quoteRepo.GetRandom(ctx)
}

func (u *QuoteUseCaseImpl) GetQuotesByFilter(ctx context.Context, author *string, content *string) ([]*entities.Quote, error) {
	filter := &repositories.QuoteFilter{
		Author:  author,
		Content: content,
	}
	return u.quoteRepo.GetByFilter(ctx, filter)
}

func (u *QuoteUseCaseImpl) UpdateQuote(ctx context.Context, id string, content, author string) error {
	quoteID, err := value_objects.NewQuoteIDFromString(id)
	if err != nil {
		return err
	}

	// Get existing quote
	existingQuote, err := u.quoteRepo.GetByID(ctx, quoteID)
	if err != nil {
		return err
	}

	// Create new quote with updated values
	updatedQuote, err := entities.NewQuote(content, author)
	if err != nil {
		return err
	}

	// Use existing ID
	updatedQuote = entities.NewQuoteFromExisting(
		existingQuote.ID(),
		updatedQuote.Content(),
		updatedQuote.Author(),
		existingQuote.CreatedAt(),
		existingQuote.UpdatedAt(),
		existingQuote.DeletedAt(),
	)

	return u.quoteRepo.Update(ctx, updatedQuote)
}

func (u *QuoteUseCaseImpl) DeleteQuote(ctx context.Context, id string) error {
	quoteID, err := value_objects.NewQuoteIDFromString(id)
	if err != nil {
		return err
	}

	return u.quoteRepo.Delete(ctx, quoteID)
}
