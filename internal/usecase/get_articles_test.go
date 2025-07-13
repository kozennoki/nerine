package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository/mocks"
	"github.com/kozennoki/nerine/internal/usecase"
	"go.uber.org/mock/gomock"
)

func TestGetArticles_Exec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     usecase.GetArticlesUsecaseInput
		setupMock func(*mocks.MockArticleRepository)
		wantLen   int
		wantPage  int
		wantTotal int
		wantErr   bool
	}{
		{
			name: "success with default page and limit",
			input: usecase.GetArticlesUsecaseInput{
				Page:  0,
				Limit: 0,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
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
				}
				m.EXPECT().CountArticles(gomock.Any()).Return(1, nil)
				m.EXPECT().GetArticles(gomock.Any(), 10, 0).Return(articles, nil)
			},
			wantLen:   1,
			wantPage:  1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "success with custom pagination",
			input: usecase.GetArticlesUsecaseInput{
				Page:  2,
				Limit: 5,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				articles := make([]*entity.Article, 5)
				for i := 0; i < 5; i++ {
					articles[i] = &entity.Article{
						ID:          string(rune('1' + i)),
						Title:       "テスト記事",
						Image:       "https://example.com/image.jpg",
						Category:    entity.Category{Slug: "tech", Name: "技術"},
						Description: "説明",
						Body:        "本文",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}
				}
				m.EXPECT().CountArticles(gomock.Any()).Return(15, nil)
				m.EXPECT().GetArticles(gomock.Any(), 5, 5).Return(articles, nil)
			},
			wantLen:   5,
			wantPage:  2,
			wantTotal: 15,
			wantErr:   false,
		},
		{
			name: "count repository error",
			input: usecase.GetArticlesUsecaseInput{
				Page:  1,
				Limit: 10,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				m.EXPECT().CountArticles(gomock.Any()).Return(0, ErrRepository)
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "get articles repository error",
			input: usecase.GetArticlesUsecaseInput{
				Page:  1,
				Limit: 10,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				m.EXPECT().CountArticles(gomock.Any()).Return(10, nil)
				m.EXPECT().GetArticles(gomock.Any(), 10, 0).Return(nil, ErrRepository)
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

			uc := usecase.NewGetArticles(mockRepo)

			got, err := uc.Exec(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticles.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got.Articles) != tt.wantLen {
					t.Errorf("GetArticles.Exec() articles length = %v, want %v", len(got.Articles), tt.wantLen)
				}
				if got.Pagination.Page != tt.wantPage {
					t.Errorf("GetArticles.Exec() page = %v, want %v", got.Pagination.Page, tt.wantPage)
				}
				if got.Pagination.Total != tt.wantTotal {
					t.Errorf("GetArticles.Exec() total = %v, want %v", got.Pagination.Total, tt.wantTotal)
				}
			}
		})
	}
}