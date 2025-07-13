package usecase_test

import (
	"context"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository/mocks"
	"github.com/kozennoki/nerine/internal/usecase"
)

func TestGetLatestArticles_Exec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       usecase.GetLatestArticlesUsecaseInput
		setupMock   func(*mocks.MockArticleRepository)
		wantLen     int
		wantErr     bool
		expectLimit int
	}{
		{
			name: "success with default limit",
			input: usecase.GetLatestArticlesUsecaseInput{
				Limit: 0,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				articles := []*entity.Article{
					{
						ID:          "1",
						Title:       "Latest Article 1",
						Image:       "https://example.com/image1.jpg",
						Category:    entity.Category{Slug: "tech", Name: "Technology"},
						Description: "Description 1",
						Body:        "Body 1",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.EXPECT().GetLatestArticles(gomock.Any(), 5).Return(articles, nil)
			},
			wantLen:     1,
			wantErr:     false,
			expectLimit: 5,
		},
		{
			name: "success with custom limit",
			input: usecase.GetLatestArticlesUsecaseInput{
				Limit: 10,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				articles := make([]*entity.Article, 10)
				for i := 0; i < 10; i++ {
					articles[i] = &entity.Article{
						ID:          string(rune('1' + i)),
						Title:       "Latest Article",
						Image:       "https://example.com/image.jpg",
						Category:    entity.Category{Slug: "tech", Name: "Technology"},
						Description: "Description",
						Body:        "Body",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}
				}
				m.EXPECT().GetLatestArticles(gomock.Any(), 10).Return(articles, nil)
			},
			wantLen:     10,
			wantErr:     false,
			expectLimit: 10,
		},
		{
			name: "limit exceeds maximum",
			input: usecase.GetLatestArticlesUsecaseInput{
				Limit: 25,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				articles := make([]*entity.Article, 20)
				for i := 0; i < 20; i++ {
					articles[i] = &entity.Article{
						ID:          string(rune('1' + i)),
						Title:       "Latest Article",
						Image:       "https://example.com/image.jpg",
						Category:    entity.Category{Slug: "tech", Name: "Technology"},
						Description: "Description",
						Body:        "Body",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}
				}
				m.EXPECT().GetLatestArticles(gomock.Any(), 20).Return(articles, nil)
			},
			wantLen:     20,
			wantErr:     false,
			expectLimit: 20,
		},
		{
			name: "repository error",
			input: usecase.GetLatestArticlesUsecaseInput{
				Limit: 5,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				m.EXPECT().GetLatestArticles(gomock.Any(), 5).Return(nil, ErrRepository)
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockArticleRepository(ctrl)
			tt.setupMock(mockRepo)

			uc := usecase.NewGetLatestArticles(mockRepo)

			got, err := uc.Exec(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestArticles.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Articles) != tt.wantLen {
				t.Errorf("GetLatestArticles.Exec() articles length = %v, want %v", len(got.Articles), tt.wantLen)
			}
		})
	}
}
