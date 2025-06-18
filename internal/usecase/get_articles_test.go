package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository/mocks"
	"github.com/kozennoki/nerine/internal/usecase"
	"go.uber.org/mock/gomock"
)

func TestGetArticles_Exec(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepository(ctrl)

	tests := []struct {
		name          string
		input         usecase.GetArticlesUsecaseInput
		setupMock     func()
		expected      usecase.GetArticlesUsecaseOutput
		expectedError error
	}{
		{
			name: "正常系: 記事一覧が正常に取得できる",
			input: usecase.GetArticlesUsecaseInput{
				Limit:  10,
				Offset: 0,
			},
			setupMock: func() {
				articles := []*entity.Article{
					{
						ID:          "1",
						Title:       "テスト記事1",
						Image:       "https://example.com/image1.jpg",
						Category:    entity.Category{Slug: "tech", Name: "技術"},
						Description: "テスト記事1の説明",
						Body:        "テスト記事1の本文",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          "2",
						Title:       "テスト記事2",
						Image:       "https://example.com/image2.jpg",
						Category:    entity.Category{Slug: "life", Name: "生活"},
						Description: "テスト記事2の説明",
						Body:        "テスト記事2の本文",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				mockRepo.EXPECT().GetArticles(gomock.Any(), 10, 0).Return(articles, nil)
			},
			expected: usecase.GetArticlesUsecaseOutput{
				Articles: []*entity.Article{
					{
						ID:          "1",
						Title:       "テスト記事1",
						Image:       "https://example.com/image1.jpg",
						Category:    entity.Category{Slug: "tech", Name: "技術"},
						Description: "テスト記事1の説明",
						Body:        "テスト記事1の本文",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          "2",
						Title:       "テスト記事2",
						Image:       "https://example.com/image2.jpg",
						Category:    entity.Category{Slug: "life", Name: "生活"},
						Description: "テスト記事2の説明",
						Body:        "テスト記事2の本文",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "正常系: 空の記事一覧",
			input: usecase.GetArticlesUsecaseInput{
				Limit:  10,
				Offset: 0,
			},
			setupMock: func() {
				mockRepo.EXPECT().GetArticles(gomock.Any(), 10, 0).Return([]*entity.Article{}, nil)
			},
			expected: usecase.GetArticlesUsecaseOutput{
				Articles: []*entity.Article{},
			},
			expectedError: nil,
		},
		{
			name: "異常系: リポジトリエラー",
			input: usecase.GetArticlesUsecaseInput{
				Limit:  10,
				Offset: 0,
			},
			setupMock: func() {
				e := errors.New("repository error")
				mockRepo.EXPECT().GetArticles(gomock.Any(), 10, 0).Return(nil, e)
			},
			expected:      usecase.GetArticlesUsecaseOutput{},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			usecase := usecase.NewGetArticles(mockRepo)
			result, err := usecase.Exec(context.Background(), tt.input)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("期待されるエラーが発生しませんでした。expected: %v", tt.expectedError)
					return
				}
				if errors.Is(err, tt.expectedError) {
					t.Errorf("エラーメッセージが一致しません。expected: %v, got: %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if len(result.Articles) != len(tt.expected.Articles) {
				t.Errorf("記事数が一致しません。expected: %d, got: %d", len(tt.expected.Articles), len(result.Articles))
				return
			}

			for i, article := range result.Articles {
				expected := tt.expected.Articles[i]
				if article.ID != expected.ID {
					t.Errorf("記事ID[%d]が一致しません。expected: %s, got: %s", i, expected.ID, article.ID)
				}
				if article.Title != expected.Title {
					t.Errorf("記事タイトル[%d]が一致しません。expected: %s, got: %s", i, expected.Title, article.Title)
				}
				if article.Category.Slug != expected.Category.Slug {
					t.Errorf("カテゴリSlug[%d]が一致しません。expected: %s, got: %s", i, expected.Category.Slug, article.Category.Slug)
				}
			}
		})
	}
}
