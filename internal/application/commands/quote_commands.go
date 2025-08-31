package commands

import (
	"errors"
)

type CreateQuoteCommand struct {
	Content string
	Author  string
}

func NewCreateQuoteCommand(content string, author string) (CreateQuoteCommand, error) {
	if content == "" {
		return CreateQuoteCommand{}, errors.New("content is required")
	}

	if author == "" {
		return CreateQuoteCommand{}, errors.New("author is required")
	}

	return CreateQuoteCommand{
		Content: content,
		Author:  author,
	}, nil
}

type UpdateQuoteCommand struct {
	ID      string
	Content string
	Author  string
}

func NewUpdateQuoteCommand(id string, content string, author string) (UpdateQuoteCommand, error) {
	if id == "" {
		return UpdateQuoteCommand{}, errors.New("id is required")
	}

	if content == "" {
		return UpdateQuoteCommand{}, errors.New("content is required")
	}

	if author == "" {
		return UpdateQuoteCommand{}, errors.New("author is required")
	}

	return UpdateQuoteCommand{
		ID:      id,
		Content: content,
		Author:  author,
	}, nil
}

type DeleteQuoteCommand struct {
	ID string
}

func NewDeleteQuoteCommand(id string) (DeleteQuoteCommand, error) {
	if id == "" {
		return DeleteQuoteCommand{}, errors.New("id is required")
	}

	return DeleteQuoteCommand{
		ID: id,
	}, nil
}
