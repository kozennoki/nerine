package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/kozennoki/nerine/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func TestArticleHandler_GetArticles(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetArticlesUsecase := mocks.NewMockGetArticlesUsecase(ctrl)
	mockGetArticleByIDUsecase := mocks.NewMockGetArticleByIDUsecase(ctrl)

	handler := NewArticleHandler(mockGetArticlesUsecase, mockGetArticleByIDUsecase)

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "正常系: デフォルトパラメータで記事一覧取得",
			queryParams: "",
			setupMock: func() {
				articles := []*entity.Article{
					{
						ID:          "1",
						Title:       "テスト記事1",
						Image:       "https://example.com/image1.jpg",
						Category:    entity.Category{Slug: "tech", Name: "技術"},
						Description: "テスト記事1の説明",
						Body:        "テスト記事1の本文",
						CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				output := usecase.GetArticlesUsecaseOutput{Articles: articles}
				mockGetArticlesUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetArticlesUsecaseInput{Limit: 10, Offset: 0}).
					Return(output, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"articles":[{"ID":"1","Title":"テスト記事1","Image":"https://example.com/image1.jpg","Category":{"Slug":"tech","Name":"技術"},"Description":"テスト記事1の説明","Body":"テスト記事1の本文","CreatedAt":"2025-01-01T00:00:00Z","UpdatedAt":"2025-01-01T00:00:00Z"}]}`,
		},
		{
			name:        "正常系: クエリパラメータ指定での記事一覧取得",
			queryParams: "?limit=5&offset=10",
			setupMock: func() {
				output := usecase.GetArticlesUsecaseOutput{Articles: []*entity.Article{}}
				mockGetArticlesUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetArticlesUsecaseInput{Limit: 5, Offset: 10}).
					Return(output, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"articles":[]}`,
		},
		{
			name:        "正常系: 無効なクエリパラメータはデフォルト値使用",
			queryParams: "?limit=invalid&offset=-1",
			setupMock: func() {
				output := usecase.GetArticlesUsecaseOutput{Articles: []*entity.Article{}}
				mockGetArticlesUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetArticlesUsecaseInput{Limit: 10, Offset: 0}).
					Return(output, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"articles":[]}`,
		},
		{
			name:        "異常系: Usecaseエラー",
			queryParams: "",
			setupMock: func() {
				mockGetArticlesUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetArticlesUsecaseInput{Limit: 10, Offset: 0}).
					Return(usecase.GetArticlesUsecaseOutput{}, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"detail":"usecase error","error":"Failed to get articles"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/articles"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetArticles(c)

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコードが一致しません。expected: %d, got: %d", tt.expectedStatus, rec.Code)
			}

			body := strings.TrimSpace(rec.Body.String())
			if body != tt.expectedBody {
				t.Errorf("レスポンスボディが一致しません。\nexpected: %s\ngot: %s", tt.expectedBody, body)
			}
		})
	}
}

func TestArticleHandler_GetArticleByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetArticlesUsecase := mocks.NewMockGetArticlesUsecase(ctrl)
	mockGetArticleByIDUsecase := mocks.NewMockGetArticleByIDUsecase(ctrl)

	handler := NewArticleHandler(mockGetArticlesUsecase, mockGetArticleByIDUsecase)

	tests := []struct {
		name           string
		articleID      string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "正常系: 記事詳細取得",
			articleID: "test-article-1",
			setupMock: func() {
				article := &entity.Article{
					ID:          "test-article-1",
					Title:       "テスト記事1",
					Image:       "https://example.com/image1.jpg",
					Category:    entity.Category{Slug: "tech", Name: "技術"},
					Description: "テスト記事1の説明",
					Body:        "テスト記事1の本文",
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				output := usecase.GetArticleByIDUsecaseOutput{Article: article}
				mockGetArticleByIDUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetArticleByIDUsecaseInput{ID: "test-article-1"}).
					Return(output, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"article":{"ID":"test-article-1","Title":"テスト記事1","Image":"https://example.com/image1.jpg","Category":{"Slug":"tech","Name":"技術"},"Description":"テスト記事1の説明","Body":"テスト記事1の本文","CreatedAt":"2025-01-01T00:00:00Z","UpdatedAt":"2025-01-01T00:00:00Z"}}`,
		},
		{
			name:           "異常系: 記事IDが空",
			articleID:      "",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Article ID is required"}`,
		},
		{
			name:      "異常系: Usecaseエラー",
			articleID: "non-existent-id",
			setupMock: func() {
				mockGetArticleByIDUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetArticleByIDUsecaseInput{ID: "non-existent-id"}).
					Return(usecase.GetArticleByIDUsecaseOutput{}, errors.New("article not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"detail":"article not found","error":"Failed to get article"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/articles/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.articleID)

			err := handler.GetArticleByID(c)

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコードが一致しません。expected: %d, got: %d", tt.expectedStatus, rec.Code)
			}

			body := strings.TrimSpace(rec.Body.String())
			if body != tt.expectedBody {
				t.Errorf("レスポンスボディが一致しません。\nexpected: %s\ngot: %s", tt.expectedBody, body)
			}
		})
	}
}
