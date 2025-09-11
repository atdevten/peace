package entities

import (
	"time"

	"github.com/atdevten/peace/internal/domain/value_objects"
)

type Tag struct {
	id          *value_objects.TagID
	name        *value_objects.TagName
	description string
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
}

func NewTag(
	name string,
	description string,
) (*Tag, error) {
	nameVO, err := value_objects.NewTagName(name)
	if err != nil {
		return nil, err
	}

	return &Tag{
		id:          value_objects.NewTagID(),
		name:        nameVO,
		description: description,
	}, nil
}

func NewTagFromExisting(
	id *value_objects.TagID,
	name *value_objects.TagName,
	description string,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *Tag {
	return &Tag{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		deletedAt:   deletedAt,
	}
}

func (t *Tag) ID() *value_objects.TagID {
	return t.id
}

func (t *Tag) Name() *value_objects.TagName {
	return t.name
}

func (t *Tag) Description() string {
	return t.description
}

func (t *Tag) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Tag) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Tag) DeletedAt() *time.Time {
	return t.deletedAt
}

func (t *Tag) UpdateName(name string) error {
	nameVO, err := value_objects.NewTagName(name)
	if err != nil {
		return err
	}
	t.name = nameVO
	return nil
}

func (t *Tag) UpdateDescription(description string) {
	t.description = description
}
