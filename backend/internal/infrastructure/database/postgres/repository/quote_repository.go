package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
	"github.com/atdevten/peace/internal/infrastructure/database/postgres/models"
	"gorm.io/gorm"
)

type QuoteRepository struct {
	db *gorm.DB
}

func NewPostgreSQLQuoteRepository(db *gorm.DB) repositories.QuoteRepository {
	return &QuoteRepository{
		db: db,
	}
}

func (r *QuoteRepository) Create(ctx context.Context, quote *entities.Quote) error {
	model := models.Quote{
		Content: quote.Content().Value(),
		Author:  quote.Author().Value(),
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("failed to create quote: %w", err)
	}

	return nil
}

func (r *QuoteRepository) GetByID(ctx context.Context, id *value_objects.QuoteID) (*entities.Quote, error) {
	var model models.Quote

	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id.Value()).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("quote not found")
		}
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	return r.mapToEntity(&model)
}

func (r *QuoteRepository) GetByFilter(ctx context.Context, filter *repositories.QuoteFilter) ([]*entities.Quote, error) {
	var models []models.Quote
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	if filter.ID != nil {
		query = query.Where("id = ?", filter.ID.Value())
	}
	if filter.Author != nil {
		query = query.Where("author ILIKE ?", "%"+*filter.Author+"%")
	}
	if filter.Content != nil {
		query = query.Where("content ILIKE ?", "%"+*filter.Content+"%")
	}

	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get quotes: %w", err)
	}

	var quotes []*entities.Quote
	for _, model := range models {
		quote, err := r.mapToEntity(&model)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (r *QuoteRepository) GetAll(ctx context.Context) ([]*entities.Quote, error) {
	var models []models.Quote

	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get quotes: %w", err)
	}

	var quotes []*entities.Quote
	for _, model := range models {
		quote, err := r.mapToEntity(&model)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (r *QuoteRepository) GetRandom(ctx context.Context) (*entities.Quote, error) {
	var model models.Quote

	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("RANDOM()").First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no quotes found")
		}
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return r.mapToEntity(&model)
}

func (r *QuoteRepository) Update(ctx context.Context, quote *entities.Quote) error {
	result := r.db.WithContext(ctx).Model(&models.Quote{}).
		Where("id = ? AND deleted_at IS NULL", quote.ID().Value()).
		Updates(map[string]interface{}{
			"content":    quote.Content().Value(),
			"author":     quote.Author().Value(),
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update quote: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("quote not found")
	}

	return nil
}

func (r *QuoteRepository) Delete(ctx context.Context, id *value_objects.QuoteID) error {
	result := r.db.WithContext(ctx).Model(&models.Quote{}).
		Where("id = ? AND deleted_at IS NULL", id.Value()).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to delete quote: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("quote not found")
	}

	return nil
}

func (r *QuoteRepository) mapToEntity(model *models.Quote) (*entities.Quote, error) {
	id := value_objects.NewQuoteIDFromInt(model.ID)
	content, err := value_objects.NewContent(model.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create content value object: %w", err)
	}

	author, err := value_objects.NewAuthor(model.Author)
	if err != nil {
		return nil, fmt.Errorf("failed to create author value object: %w", err)
	}

	return entities.NewQuoteFromExisting(
		id,
		content,
		author,
		model.CreatedAt,
		model.UpdatedAt,
		model.DeletedAt,
	), nil
}
