package usecases

import (
	"context"

	"github.com/atdevten/peace/internal/application/commands"
	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/domain/value_objects"
)

type TagUseCase interface {
	CreateTag(ctx context.Context, cmd *commands.CreateTagCommand) (*entities.Tag, error)
	GetTagByID(ctx context.Context, id int) (*entities.Tag, error)
	GetTagByName(ctx context.Context, name string) (*entities.Tag, error)
	GetAllTags(ctx context.Context) ([]*entities.Tag, error)
	UpdateTag(ctx context.Context, id int, cmd *commands.UpdateTagCommand) (*entities.Tag, error)
	DeleteTag(ctx context.Context, id int) error
	GetTagsByQuoteID(ctx context.Context, quoteID int) ([]*entities.Tag, error)
	AddTagToQuote(ctx context.Context, quoteID int, cmd *commands.AddTagToQuoteCommand) error
	RemoveTagFromQuote(ctx context.Context, quoteID int, cmd *commands.RemoveTagFromQuoteCommand) error
}

type tagUseCase struct {
	tagRepository   repositories.TagRepository
	quoteRepository repositories.QuoteRepository
}

func NewTagUseCase(
	tagRepository repositories.TagRepository,
	quoteRepository repositories.QuoteRepository,
) TagUseCase {
	return &tagUseCase{
		tagRepository:   tagRepository,
		quoteRepository: quoteRepository,
	}
}

func (uc *tagUseCase) CreateTag(ctx context.Context, cmd *commands.CreateTagCommand) (*entities.Tag, error) {
	tag, err := entities.NewTag(cmd.Name, cmd.Description)
	if err != nil {
		return nil, err
	}

	err = uc.tagRepository.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (uc *tagUseCase) GetTagByID(ctx context.Context, id int) (*entities.Tag, error) {
	tagID := value_objects.NewTagIDFromInt(id)
	return uc.tagRepository.GetByID(ctx, tagID)
}

func (uc *tagUseCase) GetTagByName(ctx context.Context, name string) (*entities.Tag, error) {
	filter := &repositories.TagFilter{
		Name: &name,
	}
	return uc.tagRepository.GetByFilter(ctx, filter)
}

func (uc *tagUseCase) GetAllTags(ctx context.Context) ([]*entities.Tag, error) {
	return uc.tagRepository.GetAll(ctx)
}

func (uc *tagUseCase) UpdateTag(ctx context.Context, id int, cmd *commands.UpdateTagCommand) (*entities.Tag, error) {
	tagID := value_objects.NewTagIDFromInt(id)

	tag, err := uc.tagRepository.GetByID(ctx, tagID)
	if err != nil {
		return nil, err
	}

	err = tag.UpdateName(cmd.Name)
	if err != nil {
		return nil, err
	}

	tag.UpdateDescription(cmd.Description)

	err = uc.tagRepository.Update(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (uc *tagUseCase) DeleteTag(ctx context.Context, id int) error {
	tagID := value_objects.NewTagIDFromInt(id)
	return uc.tagRepository.Delete(ctx, tagID)
}

func (uc *tagUseCase) GetTagsByQuoteID(ctx context.Context, quoteID int) ([]*entities.Tag, error) {
	quoteIDVO := value_objects.NewQuoteIDFromInt(quoteID)
	return uc.tagRepository.GetByQuoteID(ctx, quoteIDVO)
}

func (uc *tagUseCase) AddTagToQuote(ctx context.Context, quoteID int, cmd *commands.AddTagToQuoteCommand) error {
	// Verify quote exists
	quoteIDVO := value_objects.NewQuoteIDFromInt(quoteID)
	_, err := uc.quoteRepository.GetByID(ctx, quoteIDVO)
	if err != nil {
		return err
	}

	// Verify tag exists
	tagIDVO := value_objects.NewTagIDFromInt(cmd.TagID)
	_, err = uc.tagRepository.GetByID(ctx, tagIDVO)
	if err != nil {
		return err
	}

	return uc.tagRepository.AddTagToQuote(ctx, quoteIDVO, tagIDVO)
}

func (uc *tagUseCase) RemoveTagFromQuote(ctx context.Context, quoteID int, cmd *commands.RemoveTagFromQuoteCommand) error {
	quoteIDVO := value_objects.NewQuoteIDFromInt(quoteID)
	tagIDVO := value_objects.NewTagIDFromInt(cmd.TagID)

	return uc.tagRepository.RemoveTagFromQuote(ctx, quoteIDVO, tagIDVO)
}
