package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository/mocks"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetZennArticles_Exec_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepository(ctrl)
	useCase := usecase.NewGetZennArticles(mockRepo)

	expectedArticles := []*entity.Article{
		{
			ID:    "123",
			Title: "📝Test Zenn Article",
			Category: entity.Category{
				Slug: "zenn",
				Name: "Zenn",
			},
			Description: "Zenn記事 - test-article",
			Body:        "",
			CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			PublishedAt: time.Now(),
		},
	}

	input := usecase.GetZennArticlesUsecaseInput{
		Page:  1,
		Limit: 10,
	}

	mockRepo.EXPECT().
		GetArticles(gomock.Any(), 10, 0).
		Return(expectedArticles, nil).
		Times(1)

	output, err := useCase.Exec(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, expectedArticles, output.Articles)
	assert.Equal(t, 1, output.Pagination.Page)
	assert.Equal(t, 10, output.Pagination.Limit)
	assert.Equal(t, 0, output.Pagination.Total)
	assert.Equal(t, 0, output.Pagination.TotalPages)
}

func TestGetZennArticles_Exec_PaginationCalculation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		page           int
		limit          int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "First page",
			page:           1,
			limit:          10,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Second page",
			page:           2,
			limit:          10,
			expectedLimit:  10,
			expectedOffset: 10,
		},
		{
			name:           "Default limit when zero",
			page:           1,
			limit:          0,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Large limit capped at maximum",
			page:           1,
			limit:          200,
			expectedLimit:  100,
			expectedOffset: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockArticleRepository(ctrl)
			useCase := usecase.NewGetZennArticles(mockRepo)

			input := usecase.GetZennArticlesUsecaseInput{
				Page:  tc.page,
				Limit: tc.limit,
			}

			mockRepo.EXPECT().
				GetArticles(gomock.Any(), tc.expectedLimit, tc.expectedOffset).
				Return([]*entity.Article{}, nil).
				Times(1)

			_, err := useCase.Exec(context.Background(), input)
			assert.NoError(t, err)
		})
	}
}

func TestGetZennArticles_Exec_RepositoryError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepository(ctrl)
	useCase := usecase.NewGetZennArticles(mockRepo)

	expectedError := errors.New("repository error")

	input := usecase.GetZennArticlesUsecaseInput{
		Page:  1,
		Limit: 10,
	}

	mockRepo.EXPECT().
		GetArticles(gomock.Any(), 10, 0).
		Return(nil, expectedError).
		Times(1)

	output, err := useCase.Exec(context.Background(), input)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, output.Articles)
}

func TestGetZennArticles_Exec_EmptyResponse(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepository(ctrl)
	useCase := usecase.NewGetZennArticles(mockRepo)

	input := usecase.GetZennArticlesUsecaseInput{
		Page:  1,
		Limit: 10,
	}

	mockRepo.EXPECT().
		GetArticles(gomock.Any(), 10, 0).
		Return([]*entity.Article{}, nil).
		Times(1)

	output, err := useCase.Exec(context.Background(), input)

	require.NoError(t, err)
	assert.Empty(t, output.Articles)
	assert.Equal(t, 1, output.Pagination.Page)
	assert.Equal(t, 10, output.Pagination.Limit)
}

func TestGetZennArticles_Exec_ContextCancellation(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepository(ctrl)
	useCase := usecase.NewGetZennArticles(mockRepo)

	input := usecase.GetZennArticlesUsecaseInput{
		Page:  1,
		Limit: 10,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	expectedError := context.Canceled

	mockRepo.EXPECT().
		GetArticles(gomock.Any(), 10, 0).
		Return(nil, expectedError).
		Times(1)

	output, err := useCase.Exec(ctx, input)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, output.Articles)
}
