package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/infrastructure/utils"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/kozennoki/nerine/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

// testArticleHandlerMocks holds all mocks for ArticleHandler testing
type testArticleHandlerMocks struct {
	getArticlesUsecase           *mocks.MockGetArticlesUsecase
	getArticleByIDUsecase        *mocks.MockGetArticleByIDUsecase
	getPopularArticlesUsecase    *mocks.MockGetPopularArticlesUsecase
	getLatestArticlesUsecase     *mocks.MockGetLatestArticlesUsecase
	getArticlesByCategoryUsecase *mocks.MockGetArticlesByCategoryUsecase
}

// createTestArticleHandler creates ArticleHandler with mocks for testing
func createTestArticleHandler(ctrl *gomock.Controller) (*ArticleHandler, *testArticleHandlerMocks) {
	mocks := &testArticleHandlerMocks{
		getArticlesUsecase:           mocks.NewMockGetArticlesUsecase(ctrl),
		getArticleByIDUsecase:        mocks.NewMockGetArticleByIDUsecase(ctrl),
		getPopularArticlesUsecase:    mocks.NewMockGetPopularArticlesUsecase(ctrl),
		getLatestArticlesUsecase:     mocks.NewMockGetLatestArticlesUsecase(ctrl),
		getArticlesByCategoryUsecase: mocks.NewMockGetArticlesByCategoryUsecase(ctrl),
	}

	handler := NewArticleHandler(
		mocks.getArticlesUsecase,
		mocks.getArticleByIDUsecase,
		mocks.getPopularArticlesUsecase,
		mocks.getLatestArticlesUsecase,
		mocks.getArticlesByCategoryUsecase,
	)

	return handler, mocks
}

func TestArticleHandler_GetArticles(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, testMocks := createTestArticleHandler(ctrl)

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
				pagination := utils.Pagination{
					Total:      1,
					Page:       1,
					Limit:      10,
					TotalPages: 1,
				}
				expectedInput := usecase.GetArticlesUsecaseInput{
					Page:  1,
					Limit: 10,
				}
				expectedOutput := usecase.GetArticlesUsecaseOutput{
					Articles:   articles,
					Pagination: pagination,
				}
				testMocks.getArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"articles"`,
		},
		{
			name:        "正常系: カスタムパラメータで記事一覧取得",
			queryParams: "?page=2&limit=5",
			setupMock: func() {
				articles := []*entity.Article{}
				pagination := utils.Pagination{
					Total:      10,
					Page:       2,
					Limit:      5,
					TotalPages: 2,
				}
				expectedInput := usecase.GetArticlesUsecaseInput{
					Page:  2,
					Limit: 5,
				}
				expectedOutput := usecase.GetArticlesUsecaseOutput{
					Articles:   articles,
					Pagination: pagination,
				}
				testMocks.getArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"pagination"`,
		},
		{
			name:        "異常系: Usecaseエラー",
			queryParams: "",
			setupMock: func() {
				expectedInput := usecase.GetArticlesUsecaseInput{
					Page:  1,
					Limit: 10,
				}
				testMocks.getArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(usecase.GetArticlesUsecaseOutput{}, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error"`,
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
				t.Errorf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %s, got %s", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestArticleHandler_GetArticleByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, testMocks := createTestArticleHandler(ctrl)

	tests := []struct {
		name           string
		articleID      string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "正常系: 記事詳細取得",
			articleID: "123",
			setupMock: func() {
				article := &entity.Article{
					ID:          "123",
					Title:       "テスト記事",
					Image:       "https://example.com/image.jpg",
					Category:    entity.Category{Slug: "tech", Name: "技術"},
					Description: "テスト記事の説明",
					Body:        "テスト記事の本文",
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				}
				expectedInput := usecase.GetArticleByIDUsecaseInput{ID: "123"}
				expectedOutput := usecase.GetArticleByIDUsecaseOutput{Article: article}
				testMocks.getArticleByIDUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"article"`,
		},
		{
			name:      "異常系: 記事IDが空",
			articleID: "",
			setupMock: func() {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:      "異常系: Usecaseエラー",
			articleID: "123",
			setupMock: func() {
				expectedInput := usecase.GetArticleByIDUsecaseInput{ID: "123"}
				testMocks.getArticleByIDUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(usecase.GetArticleByIDUsecaseOutput{}, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/articles/"+tt.articleID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.articleID)

			err := handler.GetArticleByID(c)

			if err != nil {
				t.Errorf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %s, got %s", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestArticleHandler_GetPopularArticles(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, testMocks := createTestArticleHandler(ctrl)

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "正常系: デフォルトlimitで人気記事取得",
			queryParams: "",
			setupMock: func() {
				articles := []*entity.Article{
					{
						ID:          "1",
						Title:       "人気記事1",
						Image:       "https://example.com/popular1.jpg",
						Category:    entity.Category{Slug: "tech", Name: "技術"},
						Description: "人気記事1の説明",
						Body:        "人気記事1の本文",
						CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				expectedInput := usecase.GetPopularArticlesUsecaseInput{Limit: 5}
				expectedOutput := usecase.GetPopularArticlesUsecaseOutput{Articles: articles}
				testMocks.getPopularArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"articles"`,
		},
		{
			name:        "正常系: カスタムlimitで人気記事取得",
			queryParams: "?limit=10",
			setupMock: func() {
				articles := []*entity.Article{}
				expectedInput := usecase.GetPopularArticlesUsecaseInput{Limit: 10}
				expectedOutput := usecase.GetPopularArticlesUsecaseOutput{Articles: articles}
				testMocks.getPopularArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"articles"`,
		},
		{
			name:        "異常系: Usecaseエラー",
			queryParams: "",
			setupMock: func() {
				expectedInput := usecase.GetPopularArticlesUsecaseInput{Limit: 5}
				testMocks.getPopularArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(usecase.GetPopularArticlesUsecaseOutput{}, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/articles/popular"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetPopularArticles(c)

			if err != nil {
				t.Errorf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %s, got %s", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestArticleHandler_GetLatestArticles(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, testMocks := createTestArticleHandler(ctrl)

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "正常系: デフォルトlimitで最新記事取得",
			queryParams: "",
			setupMock: func() {
				articles := []*entity.Article{
					{
						ID:          "1",
						Title:       "最新記事1",
						Image:       "https://example.com/latest1.jpg",
						Category:    entity.Category{Slug: "news", Name: "ニュース"},
						Description: "最新記事1の説明",
						Body:        "最新記事1の本文",
						CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					},
				}
				expectedInput := usecase.GetLatestArticlesUsecaseInput{Limit: 5}
				expectedOutput := usecase.GetLatestArticlesUsecaseOutput{Articles: articles}
				testMocks.getLatestArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"articles"`,
		},
		{
			name:        "正常系: カスタムlimitで最新記事取得",
			queryParams: "?limit=15",
			setupMock: func() {
				articles := []*entity.Article{}
				expectedInput := usecase.GetLatestArticlesUsecaseInput{Limit: 15}
				expectedOutput := usecase.GetLatestArticlesUsecaseOutput{Articles: articles}
				testMocks.getLatestArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"articles"`,
		},
		{
			name:        "異常系: Usecaseエラー",
			queryParams: "",
			setupMock: func() {
				expectedInput := usecase.GetLatestArticlesUsecaseInput{Limit: 5}
				testMocks.getLatestArticlesUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(usecase.GetLatestArticlesUsecaseOutput{}, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/articles/latest"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetLatestArticles(c)

			if err != nil {
				t.Errorf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %s, got %s", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestArticleHandler_GetArticlesByCategory(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, testMocks := createTestArticleHandler(ctrl)

	tests := []struct {
		name           string
		categorySlug   string
		queryParams    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:         "正常系: デフォルトパラメータでカテゴリ別記事取得",
			categorySlug: "technology",
			queryParams:  "",
			setupMock: func() {
				articles := []*entity.Article{
					{
						ID:          "1",
						Title:       "技術記事1",
						Image:       "https://example.com/tech1.jpg",
						Category:    entity.Category{Slug: "technology", Name: "技術"},
						Description: "技術記事1の説明",
						Body:        "技術記事1の本文",
						CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				pagination := utils.Pagination{
					Total:      1,
					Page:       1,
					Limit:      10,
					TotalPages: 1,
				}
				expectedInput := usecase.GetArticlesByCategoryUsecaseInput{
					CategorySlug: "technology",
					Page:         1,
					Limit:        10,
				}
				expectedOutput := usecase.GetArticlesByCategoryUsecaseOutput{
					Articles:   articles,
					Pagination: pagination,
				}
				testMocks.getArticlesByCategoryUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"articles"`,
		},
		{
			name:         "正常系: カスタムパラメータでカテゴリ別記事取得",
			categorySlug: "lifestyle",
			queryParams:  "?page=2&limit=5",
			setupMock: func() {
				articles := []*entity.Article{}
				pagination := utils.Pagination{
					Total:      15,
					Page:       2,
					Limit:      5,
					TotalPages: 3,
				}
				expectedInput := usecase.GetArticlesByCategoryUsecaseInput{
					CategorySlug: "lifestyle",
					Page:         2,
					Limit:        5,
				}
				expectedOutput := usecase.GetArticlesByCategoryUsecaseOutput{
					Articles:   articles,
					Pagination: pagination,
				}
				testMocks.getArticlesByCategoryUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(expectedOutput, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"pagination"`,
		},
		{
			name:         "異常系: カテゴリスラッグが空",
			categorySlug: "",
			queryParams:  "",
			setupMock: func() {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:         "異常系: Usecaseエラー",
			categorySlug: "technology",
			queryParams:  "",
			setupMock: func() {
				expectedInput := usecase.GetArticlesByCategoryUsecaseInput{
					CategorySlug: "technology",
					Page:         1,
					Limit:        10,
				}
				testMocks.getArticlesByCategoryUsecase.EXPECT().Exec(gomock.Any(), expectedInput).Return(usecase.GetArticlesByCategoryUsecaseOutput{}, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/categories/"+tt.categorySlug+"/articles"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("slug")
			c.SetParamValues(tt.categorySlug)

			err := handler.GetArticlesByCategory(c)

			if err != nil {
				t.Errorf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %s, got %s", tt.expectedBody, rec.Body.String())
			}
		})
	}
}
