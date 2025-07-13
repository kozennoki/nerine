package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/infrastructure/utils"
	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func TestAPIHandler_GetArticles(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	tests := []struct {
		name           string
		page           *int
		limit          *int
		expectedInput  usecase.GetArticlesUsecaseInput
		mockOutput     usecase.GetArticlesUsecaseOutput
		mockError      error
		expectedStatus int
	}{
		{
			name:          "Success with default parameters",
			page:          nil,
			limit:         nil,
			expectedInput: usecase.GetArticlesUsecaseInput{Page: 1, Limit: 10},
			mockOutput: usecase.GetArticlesUsecaseOutput{
				Articles: []*entity.Article{
					{ID: "1", Title: "Test Article", Image: "test.jpg", Category: entity.Category{Slug: "tech", Name: "Technology"}},
				},
				Pagination: utils.Pagination{Total: 1, Page: 1, Limit: 10, TotalPages: 1},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:          "Success with custom parameters",
			page:          IntPtr(2),
			limit:         IntPtr(5),
			expectedInput: usecase.GetArticlesUsecaseInput{Page: 2, Limit: 5},
			mockOutput: usecase.GetArticlesUsecaseOutput{
				Articles:   []*entity.Article{},
				Pagination: utils.Pagination{Total: 0, Page: 2, Limit: 5, TotalPages: 0},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error from usecase",
			page:           nil,
			limit:          nil,
			expectedInput:  usecase.GetArticlesUsecaseInput{Page: 1, Limit: 10},
			mockOutput:     usecase.GetArticlesUsecaseOutput{},
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/articles", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			params := openapi.GetArticlesParams{
				Page:  tt.page,
				Limit: tt.limit,
			}

			mocks.GetArticlesUsecase.EXPECT().
				Exec(gomock.Any(), tt.expectedInput).
				Return(tt.mockOutput, tt.mockError)

			err := handler.GetArticles(c, params)

			if err != nil {
				t.Errorf("GetArticles() error = %v", err)
			}
			if rec.Code != tt.expectedStatus {
				t.Errorf("GetArticles() status = %v, want %v", rec.Code, tt.expectedStatus)
			}
		})
	}
}

func TestAPIHandler_GetArticleById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	tests := []struct {
		name           string
		id             string
		mockOutput     usecase.GetArticleByIDUsecaseOutput
		mockError      error
		expectedStatus int
	}{
		{
			name: "Success",
			id:   "test-id",
			mockOutput: usecase.GetArticleByIDUsecaseOutput{
				Article: &entity.Article{
					ID:    "test-id",
					Title: "Test Article",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty ID",
			id:             "",
			mockOutput:     usecase.GetArticleByIDUsecaseOutput{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error from usecase",
			id:             "test-id",
			mockOutput:     usecase.GetArticleByIDUsecaseOutput{},
			mockError:      errors.New("not found"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+tt.id, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.id != "" && tt.expectedStatus != http.StatusBadRequest {
				mocks.GetArticleByIDUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetArticleByIDUsecaseInput{ID: tt.id}).
					Return(tt.mockOutput, tt.mockError)
			}

			err := handler.GetArticleById(c, tt.id)

			if err != nil {
				t.Errorf("GetArticleById() error = %v", err)
			}
			if rec.Code != tt.expectedStatus {
				t.Errorf("GetArticleById() status = %v, want %v", rec.Code, tt.expectedStatus)
			}
		})
	}
}

func TestAPIHandler_GetPopularArticles(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	tests := []struct {
		name           string
		limit          *int
		expectedInput  usecase.GetPopularArticlesUsecaseInput
		mockOutput     usecase.GetPopularArticlesUsecaseOutput
		mockError      error
		expectedStatus int
	}{
		{
			name:          "Success with default limit",
			limit:         nil,
			expectedInput: usecase.GetPopularArticlesUsecaseInput{Limit: 5},
			mockOutput: usecase.GetPopularArticlesUsecaseOutput{
				Articles: []*entity.Article{
					{ID: "1", Title: "Popular Article"},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:          "Success with custom limit",
			limit:         IntPtr(3),
			expectedInput: usecase.GetPopularArticlesUsecaseInput{Limit: 3},
			mockOutput: usecase.GetPopularArticlesUsecaseOutput{
				Articles: []*entity.Article{},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/popular", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			params := openapi.GetPopularArticlesParams{
				Limit: tt.limit,
			}

			mocks.GetPopularArticlesUsecase.EXPECT().
				Exec(gomock.Any(), tt.expectedInput).
				Return(tt.mockOutput, tt.mockError)

			err := handler.GetPopularArticles(c, params)

			if err != nil {
				t.Errorf("GetPopularArticles() error = %v", err)
			}
			if rec.Code != tt.expectedStatus {
				t.Errorf("GetPopularArticles() status = %v, want %v", rec.Code, tt.expectedStatus)
			}
		})
	}
}

func TestAPIHandler_GetLatestArticles(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	limit := IntPtr(3)
	expectedInput := usecase.GetLatestArticlesUsecaseInput{Limit: 3}
	mockOutput := usecase.GetLatestArticlesUsecaseOutput{
		Articles: []*entity.Article{
			{ID: "1", Title: "Latest Article", CreatedAt: time.Now()},
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/latest", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	params := openapi.GetLatestArticlesParams{
		Limit: limit,
	}

	mocks.GetLatestArticlesUsecase.EXPECT().
		Exec(gomock.Any(), expectedInput).
		Return(mockOutput, nil)

	err := handler.GetLatestArticles(c, params)

	if err != nil {
		t.Errorf("GetLatestArticles() error = %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("GetLatestArticles() status = %v, want %v", rec.Code, http.StatusOK)
	}
}

func TestAPIHandler_GetArticlesByCategory(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	tests := []struct {
		name           string
		slug           string
		page           *int
		limit          *int
		expectedInput  usecase.GetArticlesByCategoryUsecaseInput
		mockOutput     usecase.GetArticlesByCategoryUsecaseOutput
		mockError      error
		expectedStatus int
	}{
		{
			name:  "Success",
			slug:  "tech",
			page:  IntPtr(1),
			limit: IntPtr(5),
			expectedInput: usecase.GetArticlesByCategoryUsecaseInput{
				CategorySlug: "tech",
				Page:         1,
				Limit:        5,
			},
			mockOutput: usecase.GetArticlesByCategoryUsecaseOutput{
				Articles: []*entity.Article{
					{ID: "1", Title: "Tech Article"},
				},
				Pagination: utils.Pagination{Total: 1, Page: 1, Limit: 5, TotalPages: 1},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty slug",
			slug:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/categories/"+tt.slug+"/articles", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			params := openapi.GetArticlesByCategoryParams{
				Page:  tt.page,
				Limit: tt.limit,
			}

			if tt.slug != "" && tt.expectedStatus != http.StatusBadRequest {
				mocks.GetArticlesByCategoryUsecase.EXPECT().
					Exec(gomock.Any(), tt.expectedInput).
					Return(tt.mockOutput, tt.mockError)
			}

			err := handler.GetArticlesByCategory(c, tt.slug, params)

			if err != nil {
				t.Errorf("GetArticlesByCategory() error = %v", err)
			}
			if rec.Code != tt.expectedStatus {
				t.Errorf("GetArticlesByCategory() status = %v, want %v", rec.Code, tt.expectedStatus)
			}
		})
	}
}
