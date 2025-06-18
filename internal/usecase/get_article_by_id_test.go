package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetArticleByID_Exec(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockArticleRepository(ctrl)

	tests := []struct {
		name          string
		input         GetArticleByIDUsecaseInput
		setupMock     func()
		expected      GetArticleByIDUsecaseOutput
		expectedError error
	}{
		{
			name: "正常系: 記事が正常に取得できる",
			input: GetArticleByIDUsecaseInput{
				ID: "test-article-1",
			},
			setupMock: func() {
				article := &entity.Article{
					ID:          "test-article-1",
					Title:       "テスト記事1",
					Image:       "https://example.com/image1.jpg",
					Category:    entity.Category{Slug: "tech", Name: "技術"},
					Description: "テスト記事1の説明",
					Body:        "テスト記事1の本文",
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				mockRepo.EXPECT().GetArticleByID(gomock.Any(), "test-article-1").Return(article, nil)
			},
			expected: GetArticleByIDUsecaseOutput{
				Article: &entity.Article{
					ID:          "test-article-1",
					Title:       "テスト記事1",
					Image:       "https://example.com/image1.jpg",
					Category:    entity.Category{Slug: "tech", Name: "技術"},
					Description: "テスト記事1の説明",
					Body:        "テスト記事1の本文",
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expectedError: nil,
		},
		{
			name: "異常系: 記事が見つからない",
			input: GetArticleByIDUsecaseInput{
				ID: "non-existent-id",
			},
			setupMock: func() {
				e := errors.New("article not found")
				mockRepo.EXPECT().GetArticleByID(gomock.Any(), "non-existent-id").Return(nil, e)
			},
			expected:      GetArticleByIDUsecaseOutput{},
			expectedError: errors.New("article not found"),
		},
		{
			name: "異常系: リポジトリエラー",
			input: GetArticleByIDUsecaseInput{
				ID: "test-article-1",
			},
			setupMock: func() {
				e := errors.New("repository error")
				mockRepo.EXPECT().GetArticleByID(gomock.Any(), "test-article-1").Return(nil, e)
			},
			expected:      GetArticleByIDUsecaseOutput{},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			usecase := NewGetArticleByID(mockRepo)
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

			if result.Article == nil {
				t.Error("記事がnilです")
				return
			}

			expected := tt.expected.Article
			if result.Article.ID != expected.ID {
				t.Errorf("記事IDが一致しません。expected: %s, got: %s", expected.ID, result.Article.ID)
			}
			if result.Article.Title != expected.Title {
				t.Errorf("記事タイトルが一致しません。expected: %s, got: %s", expected.Title, result.Article.Title)
			}
			if result.Article.Image != expected.Image {
				t.Errorf("記事画像が一致しません。expected: %s, got: %s", expected.Image, result.Article.Image)
			}
			if result.Article.Category.Slug != expected.Category.Slug {
				t.Errorf("カテゴリSlugが一致しません。expected: %s, got: %s", expected.Category.Slug, result.Article.Category.Slug)
			}
			if result.Article.Category.Name != expected.Category.Name {
				t.Errorf("カテゴリ名が一致しません。expected: %s, got: %s", expected.Category.Name, result.Article.Category.Name)
			}
			if result.Article.Description != expected.Description {
				t.Errorf("記事説明が一致しません。expected: %s, got: %s", expected.Description, result.Article.Description)
			}
			if result.Article.Body != expected.Body {
				t.Errorf("記事本文が一致しません。expected: %s, got: %s", expected.Body, result.Article.Body)
			}
		})
	}
}
