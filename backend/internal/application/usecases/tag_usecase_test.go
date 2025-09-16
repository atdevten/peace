package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/atdevten/peace/internal/application/commands"
	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/testutils/helpers"
	repositories "github.com/atdevten/peace/testutils/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTagUseCase_CreateTag(t *testing.T) {
	tests := []struct {
		name        string
		command     *commands.CreateTagCommand
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful creation",
			command: &commands.CreateTagCommand{
				Name:        "motivation",
				Description: "Motivational quotes and content",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			command: &commands.CreateTagCommand{
				Name:        "",
				Description: "Motivational quotes and content",
			},
			wantErr:     true,
			expectedErr: "tag name cannot be empty",
		},
		{
			name: "repository error",
			command: &commands.CreateTagCommand{
				Name:        "motivation",
				Description: "Motivational quotes and content",
			},
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

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			if tt.command.Name != "" {
				mockTagRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.mockError)
			}

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			tag, err := useCase.CreateTag(context.Background(), tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, tag)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tag)
			}
		})
	}
}

func TestTagUseCase_GetTagByID(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		mockTag     *entities.Tag
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful get",
			id:      123,
			mockTag: helpers.CreateTestTag(),
			wantErr: false,
		},
		{
			name:        "tag not found",
			id:          123,
			mockError:   errors.New("tag not found"),
			wantErr:     true,
			expectedErr: "tag not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			mockTagRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockTag, tt.mockError)

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			tag, err := useCase.GetTagByID(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, tag)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tag)
			}
		})
	}
}

func TestTagUseCase_GetTagByName(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		mockTag     *entities.Tag
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful get",
			tagName: "motivation",
			mockTag: helpers.CreateTestTag(),
			wantErr: false,
		},
		{
			name:        "tag not found",
			tagName:     "nonexistent",
			mockError:   errors.New("tag not found"),
			wantErr:     true,
			expectedErr: "tag not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			mockTagRepo.EXPECT().GetByFilter(gomock.Any(), gomock.Any()).Return(tt.mockTag, tt.mockError)

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			tag, err := useCase.GetTagByName(context.Background(), tt.tagName)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, tag)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tag)
			}
		})
	}
}

func TestTagUseCase_GetAllTags(t *testing.T) {
	tests := []struct {
		name        string
		mockTags    []*entities.Tag
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "successful get all",
			mockTags: []*entities.Tag{helpers.CreateTestTag()},
			wantErr:  false,
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

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			mockTagRepo.EXPECT().GetAll(gomock.Any()).Return(tt.mockTags, tt.mockError)

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			tags, err := useCase.GetAllTags(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, tags)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tags)
				assert.Len(t, tags, len(tt.mockTags))
			}
		})
	}
}

func TestTagUseCase_UpdateTag(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		command     *commands.UpdateTagCommand
		mockTag     *entities.Tag
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful update",
			id:   123,
			command: &commands.UpdateTagCommand{
				Name:        "updated-motivation",
				Description: "Updated motivational content",
			},
			mockTag: helpers.CreateTestTag(),
			wantErr: false,
		},
		{
			name: "tag not found",
			id:   123,
			command: &commands.UpdateTagCommand{
				Name:        "updated-motivation",
				Description: "Updated motivational content",
			},
			mockError:   errors.New("tag not found"),
			wantErr:     true,
			expectedErr: "tag not found",
		},
		{
			name: "empty name",
			id:   123,
			command: &commands.UpdateTagCommand{
				Name:        "",
				Description: "Updated motivational content",
			},
			mockTag:     helpers.CreateTestTag(),
			wantErr:     true,
			expectedErr: "tag name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			mockTagRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(tt.mockTag, tt.mockError)
			if tt.mockError == nil && tt.command.Name != "" {
				mockTagRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			}

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			tag, err := useCase.UpdateTag(context.Background(), tt.id, tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, tag)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tag)
			}
		})
	}
}

func TestTagUseCase_DeleteTag(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful deletion",
			id:      123,
			wantErr: false,
		},
		{
			name:        "tag not found",
			id:          123,
			mockError:   errors.New("tag not found"),
			wantErr:     true,
			expectedErr: "tag not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			mockTagRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.mockError)

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			err := useCase.DeleteTag(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTagUseCase_GetTagsByQuoteID(t *testing.T) {
	tests := []struct {
		name        string
		quoteID     int
		mockTags    []*entities.Tag
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "successful get tags by quote ID",
			quoteID:  123,
			mockTags: []*entities.Tag{helpers.CreateTestTag()},
			wantErr:  false,
		},
		{
			name:        "repository error",
			quoteID:     123,
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

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			mockTagRepo.EXPECT().GetByQuoteID(gomock.Any(), gomock.Any()).Return(tt.mockTags, tt.mockError)

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			tags, err := useCase.GetTagsByQuoteID(context.Background(), tt.quoteID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, tags)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tags)
				assert.Len(t, tags, len(tt.mockTags))
			}
		})
	}
}

func TestTagUseCase_AddTagToQuote(t *testing.T) {
	tests := []struct {
		name        string
		quoteID     int
		command     *commands.AddTagToQuoteCommand
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful add tag to quote",
			quoteID: 123,
			command: &commands.AddTagToQuoteCommand{
				TagID: 456,
			},
			wantErr: false,
		},
		{
			name:    "repository error",
			quoteID: 123,
			command: &commands.AddTagToQuoteCommand{
				TagID: 456,
			},
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

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			mockQuoteRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(helpers.CreateTestQuote(), nil)
			mockTagRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(helpers.CreateTestTag(), nil)
			mockTagRepo.EXPECT().AddTagToQuote(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.mockError)

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			err := useCase.AddTagToQuote(context.Background(), tt.quoteID, tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTagUseCase_RemoveTagFromQuote(t *testing.T) {
	tests := []struct {
		name        string
		quoteID     int
		command     *commands.RemoveTagFromQuoteCommand
		mockError   error
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "successful remove tag from quote",
			quoteID: 123,
			command: &commands.RemoveTagFromQuoteCommand{
				TagID: 456,
			},
			wantErr: false,
		},
		{
			name:    "repository error",
			quoteID: 123,
			command: &commands.RemoveTagFromQuoteCommand{
				TagID: 456,
			},
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

			// Setup mock repositories
			mockTagRepo := repositories.NewMockTagRepository(ctrl)
			mockQuoteRepo := repositories.NewMockQuoteRepository(ctrl)

			mockTagRepo.EXPECT().RemoveTagFromQuote(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.mockError)

			useCase := NewTagUseCase(mockTagRepo, mockQuoteRepo)
			err := useCase.RemoveTagFromQuote(context.Background(), tt.quoteID, tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
