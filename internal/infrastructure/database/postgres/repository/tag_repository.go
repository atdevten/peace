package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
	"github.com/atdevten/peace/internal/infrastructure/database/postgres/models"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repositories.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *entities.Tag) error {
	tagModel := &models.Tag{
		Name:        tag.Name().Value(),
		Description: tag.Description(),
	}

	result := r.db.WithContext(ctx).Create(tagModel)
	if result.Error != nil {
		return result.Error
	}

	// Update the entity with the generated ID
	tagID := value_objects.NewTagIDFromInt(tagModel.ID)
	*tag = *entities.NewTagFromExisting(
		tagID,
		tag.Name(),
		tag.Description(),
		tagModel.CreatedAt,
		tagModel.UpdatedAt,
		nil,
	)

	return nil
}

func (r *tagRepository) GetByID(ctx context.Context, id *value_objects.TagID) (*entities.Tag, error) {
	var tagModel models.Tag

	result := r.db.WithContext(ctx).Where("id = ?", id.IntValue()).First(&tagModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrTagNotFound
		}
		return nil, result.Error
	}

	tagID, tagName, err := tagModel.ToDomain()
	if err != nil {
		return nil, err
	}

	return entities.NewTagFromExisting(
		tagID,
		tagName,
		tagModel.Description,
		tagModel.CreatedAt,
		tagModel.UpdatedAt,
		&tagModel.DeletedAt.Time,
	), nil
}

func (r *tagRepository) GetByFilter(ctx context.Context, filter *repositories.TagFilter) (*entities.Tag, error) {
	var tagModel models.Tag
	var query string
	var queryValue interface{}

	if filter.ID != nil {
		query = "id = ?"
		queryValue = filter.ID.IntValue()
	} else if filter.Name != nil {
		query = "name = ?"
		queryValue = *filter.Name
	} else {
		return nil, fmt.Errorf("no search criteria provided")
	}

	result := r.db.WithContext(ctx).Where(query, queryValue).First(&tagModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrTagNotFound
		}
		return nil, result.Error
	}

	tagID, tagName, err := tagModel.ToDomain()
	if err != nil {
		return nil, err
	}

	return entities.NewTagFromExisting(
		tagID,
		tagName,
		tagModel.Description,
		tagModel.CreatedAt,
		tagModel.UpdatedAt,
		&tagModel.DeletedAt.Time,
	), nil
}

func (r *tagRepository) GetAll(ctx context.Context) ([]*entities.Tag, error) {
	var tagModels []models.Tag

	result := r.db.WithContext(ctx).Find(&tagModels)
	if result.Error != nil {
		return nil, result.Error
	}

	tags := make([]*entities.Tag, len(tagModels))
	for i, tagModel := range tagModels {
		tagID, tagName, err := tagModel.ToDomain()
		if err != nil {
			return nil, err
		}

		tags[i] = entities.NewTagFromExisting(
			tagID,
			tagName,
			tagModel.Description,
			tagModel.CreatedAt,
			tagModel.UpdatedAt,
			&tagModel.DeletedAt.Time,
		)
	}

	return tags, nil
}

func (r *tagRepository) Update(ctx context.Context, tag *entities.Tag) error {
	tagModel := &models.Tag{
		ID:          tag.ID().IntValue(),
		Name:        tag.Name().Value(),
		Description: tag.Description(),
	}

	result := r.db.WithContext(ctx).Save(tagModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *tagRepository) Delete(ctx context.Context, id *value_objects.TagID) error {
	result := r.db.WithContext(ctx).Delete(&models.Tag{}, id.IntValue())
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repositories.ErrTagNotFound
	}

	return nil
}

func (r *tagRepository) GetByQuoteID(ctx context.Context, quoteID *value_objects.QuoteID) ([]*entities.Tag, error) {
	var tagModels []models.Tag

	result := r.db.WithContext(ctx).
		Joins("JOIN quote_tags ON tags.id = quote_tags.tag_id").
		Where("quote_tags.quote_id = ?", quoteID.Value()).
		Find(&tagModels)

	if result.Error != nil {
		return nil, result.Error
	}

	tags := make([]*entities.Tag, len(tagModels))
	for i, tagModel := range tagModels {
		tagID, tagName, err := tagModel.ToDomain()
		if err != nil {
			return nil, err
		}

		tags[i] = entities.NewTagFromExisting(
			tagID,
			tagName,
			tagModel.Description,
			tagModel.CreatedAt,
			tagModel.UpdatedAt,
			&tagModel.DeletedAt.Time,
		)
	}

	return tags, nil
}

func (r *tagRepository) AddTagToQuote(ctx context.Context, quoteID *value_objects.QuoteID, tagID *value_objects.TagID) error {
	quoteTag := &models.QuoteTag{}
	quoteTag.FromDomain(quoteID, tagID)

	result := r.db.WithContext(ctx).Create(quoteTag)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *tagRepository) RemoveTagFromQuote(ctx context.Context, quoteID *value_objects.QuoteID, tagID *value_objects.TagID) error {
	result := r.db.WithContext(ctx).
		Where("quote_id = ? AND tag_id = ?", quoteID.Value(), tagID.IntValue()).
		Delete(&models.QuoteTag{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repositories.ErrTagNotFound
	}

	return nil
}
