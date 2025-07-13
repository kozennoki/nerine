package handlers

import (
	"net/http"

	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/labstack/echo/v4"
)

type APIHandler struct {
	*ArticleHandler
	*CategoryHandler
}

func NewAPIHandler(
	getArticlesUsecase usecase.GetArticlesUsecase,
	getArticleByIDUsecase usecase.GetArticleByIDUsecase,
	getPopularArticlesUsecase usecase.GetPopularArticlesUsecase,
	getLatestArticlesUsecase usecase.GetLatestArticlesUsecase,
	getArticlesByCategoryUsecase usecase.GetArticlesByCategoryUsecase,
	getCategoriesUsecase usecase.GetCategoriesUsecase,
) *APIHandler {
	return &APIHandler{
		ArticleHandler: NewArticleHandler(
			getArticlesUsecase,
			getArticleByIDUsecase,
			getPopularArticlesUsecase,
			getLatestArticlesUsecase,
			getArticlesByCategoryUsecase,
		),
		CategoryHandler: NewCategoryHandler(getCategoriesUsecase),
	}
}

func (h *APIHandler) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, openapi.HealthResponse{
		Status: "ok",
	})
}

var _ openapi.ServerInterface = (*APIHandler)(nil)