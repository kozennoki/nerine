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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAPIHandler_GetZennArticles_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

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
		},
	}

	expectedOutput := usecase.GetZennArticlesUsecaseOutput{
		Articles: expectedArticles,
		Pagination: utils.Pagination{
			Total:      1,
			Page:       1,
			Limit:      10,
			TotalPages: 1,
		},
	}

	mocks.GetZennArticlesUsecase.EXPECT().
		Exec(gomock.Any(), usecase.GetZennArticlesUsecaseInput{
			Page:  1,
			Limit: 10,
		}).
		Return(expectedOutput, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/zenn/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	params := openapi.GetZennArticlesParams{}

	err := handler.GetZennArticles(c, params)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "📝Test Zenn Article")
	assert.Contains(t, rec.Body.String(), "zenn")
}

func TestAPIHandler_GetZennArticles_WithPagination(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	expectedOutput := usecase.GetZennArticlesUsecaseOutput{
		Articles: []*entity.Article{},
		Pagination: utils.Pagination{
			Total:      50,
			Page:       2,
			Limit:      5,
			TotalPages: 10,
		},
	}

	page := 2
	limit := 5

	mocks.GetZennArticlesUsecase.EXPECT().
		Exec(gomock.Any(), usecase.GetZennArticlesUsecaseInput{
			Page:  page,
			Limit: limit,
		}).
		Return(expectedOutput, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/zenn/articles?page=2&limit=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	params := openapi.GetZennArticlesParams{
		Page:  &page,
		Limit: &limit,
	}

	err := handler.GetZennArticles(c, params)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"page":2`)
	assert.Contains(t, rec.Body.String(), `"limit":5`)
	assert.Contains(t, rec.Body.String(), `"total":50`)
	assert.Contains(t, rec.Body.String(), `"totalPages":10`)
}

func TestAPIHandler_GetZennArticles_DefaultValues(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	expectedOutput := usecase.GetZennArticlesUsecaseOutput{
		Articles: []*entity.Article{},
		Pagination: utils.Pagination{
			Total:      0,
			Page:       1,
			Limit:      10,
			TotalPages: 0,
		},
	}

	mocks.GetZennArticlesUsecase.EXPECT().
		Exec(gomock.Any(), usecase.GetZennArticlesUsecaseInput{
			Page:  1,
			Limit: 10,
		}).
		Return(expectedOutput, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/zenn/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	params := openapi.GetZennArticlesParams{}

	err := handler.GetZennArticles(c, params)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"page":1`)
	assert.Contains(t, rec.Body.String(), `"limit":10`)
}

func TestAPIHandler_GetZennArticles_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	expectedError := errors.New("usecase error")

	mocks.GetZennArticlesUsecase.EXPECT().
		Exec(gomock.Any(), usecase.GetZennArticlesUsecaseInput{
			Page:  1,
			Limit: 10,
		}).
		Return(usecase.GetZennArticlesUsecaseOutput{}, expectedError).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/zenn/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	params := openapi.GetZennArticlesParams{}

	err := handler.GetZennArticles(c, params)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to get Zenn articles")
	assert.Contains(t, rec.Body.String(), "usecase error")
}

func TestAPIHandler_GetZennArticles_EmptyResponse(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	expectedOutput := usecase.GetZennArticlesUsecaseOutput{
		Articles: []*entity.Article{},
		Pagination: utils.Pagination{
			Total:      0,
			Page:       1,
			Limit:      10,
			TotalPages: 0,
		},
	}

	mocks.GetZennArticlesUsecase.EXPECT().
		Exec(gomock.Any(), usecase.GetZennArticlesUsecaseInput{
			Page:  1,
			Limit: 10,
		}).
		Return(expectedOutput, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/zenn/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	params := openapi.GetZennArticlesParams{}

	err := handler.GetZennArticles(c, params)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"articles":[]`)
}

func TestAPIHandler_GetZennArticles_MultipleArticles(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	expectedArticles := []*entity.Article{
		{
			ID:    "123",
			Title: "📝First Article",
			Category: entity.Category{
				Slug: "zenn",
				Name: "Zenn",
			},
			Description: "Zenn記事 - first-article",
			Body:        "",
			CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:    "456",
			Title: "🚀Second Article",
			Category: entity.Category{
				Slug: "zenn",
				Name: "Zenn",
			},
			Description: "Zenn記事 - second-article",
			Body:        "",
			CreatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	expectedOutput := usecase.GetZennArticlesUsecaseOutput{
		Articles: expectedArticles,
		Pagination: utils.Pagination{
			Total:      2,
			Page:       1,
			Limit:      10,
			TotalPages: 1,
		},
	}

	mocks.GetZennArticlesUsecase.EXPECT().
		Exec(gomock.Any(), usecase.GetZennArticlesUsecaseInput{
			Page:  1,
			Limit: 10,
		}).
		Return(expectedOutput, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/zenn/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	params := openapi.GetZennArticlesParams{}

	err := handler.GetZennArticles(c, params)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "First Article")
	assert.Contains(t, rec.Body.String(), "Second Article")
	assert.Contains(t, rec.Body.String(), "📝")
	assert.Contains(t, rec.Body.String(), "🚀")
}
