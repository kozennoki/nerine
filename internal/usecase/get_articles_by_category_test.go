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

func TestGetArticlesByCategory_Exec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     usecase.GetArticlesByCategoryUsecaseInput
		setupMock func(*mocks.MockArticleRepository)
		wantLen   int
		wantPage  int
		wantTotal int
		wantErr   bool
	}{
		{
			name: "success with default page and limit",
			input: usecase.GetArticlesByCategoryUsecaseInput{
				CategorySlug: "technology",
				Page:         0,
				Limit:        0,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				articles := []*entity.Article{
					{
						ID:          "1",
						Title:       "Tech Article 1",
						Image:       "https://example.com/image1.jpg",
						Category:    entity.Category{Slug: "technology", Name: "Technology"},
						Description: "Description 1",
						Body:        "Body 1",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.EXPECT().CountArticlesByCategory(gomock.Any(), "technology").Return(1, nil)
				m.EXPECT().GetArticlesByCategory(gomock.Any(), "technology", 10, 0).Return(articles, nil)
			},
			wantLen:   1,
			wantPage:  1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "success with custom pagination",
			input: usecase.GetArticlesByCategoryUsecaseInput{
				CategorySlug: "technology",
				Page:         2,
				Limit:        5,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				articles := make([]*entity.Article, 5)
				for i := 0; i < 5; i++ {
					articles[i] = &entity.Article{
						ID:          string(rune('1' + i)),
						Title:       "Tech Article",
						Image:       "https://example.com/image.jpg",
						Category:    entity.Category{Slug: "technology", Name: "Technology"},
						Description: "Description",
						Body:        "Body",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}
				}
				m.EXPECT().CountArticlesByCategory(gomock.Any(), "technology").Return(15, nil)
				m.EXPECT().GetArticlesByCategory(gomock.Any(), "technology", 5, 5).Return(articles, nil)
			},
			wantLen:   5,
			wantPage:  2,
			wantTotal: 15,
			wantErr:   false,
		},
		{
			name: "limit exceeds maximum",
			input: usecase.GetArticlesByCategoryUsecaseInput{
				CategorySlug: "technology",
				Page:         1,
				Limit:        150,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				articles := make([]*entity.Article, 100)
				for i := 0; i < 100; i++ {
					articles[i] = &entity.Article{
						ID:          string(rune('1' + i)),
						Title:       "Tech Article",
						Image:       "https://example.com/image.jpg",
						Category:    entity.Category{Slug: "technology", Name: "Technology"},
						Description: "Description",
						Body:        "Body",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}
				}
				m.EXPECT().CountArticlesByCategory(gomock.Any(), "technology").Return(200, nil)
				m.EXPECT().GetArticlesByCategory(gomock.Any(), "technology", 100, 0).Return(articles, nil)
			},
			wantLen:   100,
			wantPage:  1,
			wantTotal: 200,
			wantErr:   false,
		},
		{
			name: "count repository error",
			input: usecase.GetArticlesByCategoryUsecaseInput{
				CategorySlug: "technology",
				Page:         1,
				Limit:        10,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				m.EXPECT().CountArticlesByCategory(gomock.Any(), "technology").Return(0, ErrRepository)
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "get articles repository error",
			input: usecase.GetArticlesByCategoryUsecaseInput{
				CategorySlug: "technology",
				Page:         1,
				Limit:        10,
			},
			setupMock: func(m *mocks.MockArticleRepository) {
				m.EXPECT().CountArticlesByCategory(gomock.Any(), "technology").Return(10, nil)
				m.EXPECT().GetArticlesByCategory(gomock.Any(), "technology", 10, 0).Return(nil, ErrRepository)
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

			uc := usecase.NewGetArticlesByCategory(mockRepo)

			got, err := uc.Exec(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticlesByCategory.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got.Articles) != tt.wantLen {
					t.Errorf("GetArticlesByCategory.Exec() articles length = %v, want %v", len(got.Articles), tt.wantLen)
				}
				if got.Pagination.Page != tt.wantPage {
					t.Errorf("GetArticlesByCategory.Exec() page = %v, want %v", got.Pagination.Page, tt.wantPage)
				}
				if got.Pagination.Total != tt.wantTotal {
					t.Errorf("GetArticlesByCategory.Exec() total = %v, want %v", got.Pagination.Total, tt.wantTotal)
				}
			}
		})
	}
}
