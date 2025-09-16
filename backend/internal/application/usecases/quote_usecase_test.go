package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/testutils/helpers"
	repositories "github.com/atdevten/peace/testutils/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestQuoteUseCaseImpl_CreateQuote(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		author      string
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful creation",
			content: "Life is beautiful",
			author:  "Anonymous",
			wantErr: false,
		},
		{
			name:        "empty content",
			content:     "",
			author:      "Anonymous",
			wantErr:     true,
			expectedErr: "content cannot be empty",
		},
		{
			name:        "empty author",
			content:     "Life is beautiful",
			author:      "",
			wantErr:     true,
			expectedErr: "author cannot be empty",
		},
		{
			name:        "repository error",
			content:     "Life is beautiful",
			author:      "Anonymous",
			mockError:   errors.New("database error"),
			wantErr:     true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockQuoteRepository(ctrl)
			if tt.content != "" && tt.author != "" {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.mockError)
			}

			useCase := NewQuoteUseCase(mockRepo)
			err := useCase.CreateQuote(context.Background(), tt.content, tt.author)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQuoteUseCaseImpl_GetQuoteByID(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		mockQuote   *entities.Quote
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:      "successful get",
			id:        "123",
			mockQuote: helpers.CreateTestQuote(),
			wantErr:   false,
		},
		{
			name:        "invalid ID",
			id:          "invalid-id",
			wantErr:     true,
			expectedErr: "invalid quote id",
		},
		{
			name:        "quote not found",
			id:          "123",
			mockError:   errors.New("quote not found"),
			wantErr:     true,
			expectedErr: "quote not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockQuoteRepository(ctrl)
			if tt.id != "invalid-id" {
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockQuote, tt.mockError)
			}

			useCase := NewQuoteUseCase(mockRepo)
			quote, err := useCase.GetQuoteByID(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, quote)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, quote)
			}
		})
	}
}

func TestQuoteUseCaseImpl_GetAllQuotes(t *testing.T) {
	tests := []struct {
		name        string
		mockQuotes  []*entities.Quote
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:       "successful get all",
			mockQuotes: []*entities.Quote{helpers.CreateTestQuote()},
			wantErr:    false,
		},
		{
			name:        "repository error",
			mockError:   errors.New("database error"),
			wantErr:     true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockQuoteRepository(ctrl)
			mockRepo.EXPECT().GetAll(gomock.Any()).Return(tt.mockQuotes, tt.mockError)

			useCase := NewQuoteUseCase(mockRepo)
			quotes, err := useCase.GetAllQuotes(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, quotes)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, quotes)
				assert.Len(t, quotes, len(tt.mockQuotes))
			}
		})
	}
}

func TestQuoteUseCaseImpl_GetRandomQuote(t *testing.T) {
	tests := []struct {
		name        string
		mockQuote   *entities.Quote
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:      "successful get random",
			mockQuote: helpers.CreateTestQuote(),
			wantErr:   false,
		},
		{
			name:        "no quotes available",
			mockError:   errors.New("no quotes found"),
			wantErr:     true,
			expectedErr: "no quotes found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockQuoteRepository(ctrl)
			mockRepo.EXPECT().GetRandom(gomock.Any()).Return(tt.mockQuote, tt.mockError)

			useCase := NewQuoteUseCase(mockRepo)
			quote, err := useCase.GetRandomQuote(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, quote)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, quote)
			}
		})
	}
}

func TestQuoteUseCaseImpl_GetQuotesByFilter(t *testing.T) {
	tests := []struct {
		name        string
		author      *string
		content     *string
		mockQuotes  []*entities.Quote
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:       "successful filter by author",
			author:     helpers.StringPtr("Anonymous"),
			mockQuotes: []*entities.Quote{helpers.CreateTestQuote()},
			wantErr:    false,
		},
		{
			name:       "successful filter by content",
			content:    helpers.StringPtr("Life"),
			mockQuotes: []*entities.Quote{helpers.CreateTestQuote()},
			wantErr:    false,
		},
		{
			name:        "repository error",
			author:      helpers.StringPtr("Anonymous"),
			mockError:   errors.New("database error"),
			wantErr:     true,
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockQuoteRepository(ctrl)
			mockRepo.EXPECT().GetByFilter(gomock.Any(), gomock.Any()).Return(tt.mockQuotes, tt.mockError)

			useCase := NewQuoteUseCase(mockRepo)
			quotes, err := useCase.GetQuotesByFilter(context.Background(), tt.author, tt.content)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, quotes)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, quotes)
				assert.Len(t, quotes, len(tt.mockQuotes))
			}
		})
	}
}

func TestQuoteUseCaseImpl_UpdateQuote(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		content     string
		author      string
		mockQuote   *entities.Quote
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:      "successful update",
			id:        "123",
			content:   "Updated content",
			author:    "Updated Author",
			mockQuote: helpers.CreateTestQuote(),
			wantErr:   false,
		},
		{
			name:        "invalid ID",
			id:          "invalid-id",
			content:     "Updated content",
			author:      "Updated Author",
			wantErr:     true,
			expectedErr: "invalid quote id",
		},
		{
			name:        "quote not found",
			id:          "123",
			content:     "Updated content",
			author:      "Updated Author",
			mockError:   errors.New("quote not found"),
			wantErr:     true,
			expectedErr: "quote not found",
		},
		{
			name:        "empty content",
			id:          "123",
			content:     "",
			author:      "Updated Author",
			mockQuote:   helpers.CreateTestQuote(),
			wantErr:     true,
			expectedErr: "content cannot be empty",
		},
		{
			name:        "empty author",
			id:          "123",
			content:     "Updated content",
			author:      "",
			mockQuote:   helpers.CreateTestQuote(),
			wantErr:     true,
			expectedErr: "author cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockQuoteRepository(ctrl)
			if tt.id != "invalid-id" {
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockQuote, tt.mockError)
				if tt.mockError == nil && tt.content != "" && tt.author != "" {
					mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				}
			}

			useCase := NewQuoteUseCase(mockRepo)
			err := useCase.UpdateQuote(context.Background(), tt.id, tt.content, tt.author)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQuoteUseCaseImpl_DeleteQuote(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		mockQuote   *entities.Quote
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:      "successful deletion",
			id:        "123",
			mockQuote: helpers.CreateTestQuote(),
			wantErr:   false,
		},
		{
			name:        "invalid ID",
			id:          "invalid-id",
			wantErr:     true,
			expectedErr: "invalid quote id",
		},
		{
			name:        "quote not found",
			id:          "123",
			mockError:   errors.New("quote not found"),
			wantErr:     true,
			expectedErr: "quote not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := repositories.NewMockQuoteRepository(ctrl)
			if tt.id != "invalid-id" {
				mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.mockError)
			}

			useCase := NewQuoteUseCase(mockRepo)
			err := useCase.DeleteQuote(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
