package handlers

import (
	"net/http"

	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/kozennoki/nerine/internal/interfaces/presenter"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/labstack/echo/v4"
)

func (h *APIHandler) GetArticles(ctx echo.Context, params openapi.GetArticlesParams) error {
	page := 1
	if params.Page != nil {
		page = *params.Page
	}

	limit := 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	input := usecase.GetArticlesUsecaseInput{
		Page:  page,
		Limit: limit,
	}

	output, err := h.getArticlesUsecase.Exec(ctx.Request().Context(), input)
	if err != nil {
		ctx.Logger().Error("Failed to get articles: ", err)
		errorMsg := presenter.ConvertErrorMessage(err)
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Error:  "Failed to get articles",
			Detail: &errorMsg,
		})
	}

	return ctx.JSON(http.StatusOK, openapi.ArticlesResponse{
		Articles:   presenter.ConvertArticles(output.Articles),
		Pagination: presenter.ConvertPagination(output.Pagination),
	})
}

func (h *APIHandler) GetArticleById(ctx echo.Context, id string) error {
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Error: "Article ID is required",
		})
	}

	input := usecase.GetArticleByIDUsecaseInput{
		ID: id,
	}

	output, err := h.getArticleByIDUsecase.Exec(ctx.Request().Context(), input)
	if err != nil {
		ctx.Logger().Error("Failed to get article by ID: ", err)
		errorMsg := presenter.ConvertErrorMessage(err)
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Error:  "Failed to get article",
			Detail: &errorMsg,
		})
	}

	return ctx.JSON(http.StatusOK, openapi.ArticleResponse{
		Article: presenter.ConvertArticle(output.Article),
	})
}

func (h *APIHandler) GetPopularArticles(ctx echo.Context, params openapi.GetPopularArticlesParams) error {
	limit := 5
	if params.Limit != nil {
		limit = *params.Limit
	}

	input := usecase.GetPopularArticlesUsecaseInput{
		Limit: limit,
	}

	output, err := h.getPopularArticlesUsecase.Exec(ctx.Request().Context(), input)
	if err != nil {
		ctx.Logger().Error("Failed to get popular articles: ", err)
		errorMsg := presenter.ConvertErrorMessage(err)
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Error:  "Failed to get articles",
			Detail: &errorMsg,
		})
	}

	return ctx.JSON(http.StatusOK, openapi.ArticlesResponse{
		Articles: presenter.ConvertArticles(output.Articles),
	})
}

func (h *APIHandler) GetLatestArticles(ctx echo.Context, params openapi.GetLatestArticlesParams) error {
	limit := 5
	if params.Limit != nil {
		limit = *params.Limit
	}

	input := usecase.GetLatestArticlesUsecaseInput{
		Limit: limit,
	}

	output, err := h.getLatestArticlesUsecase.Exec(ctx.Request().Context(), input)
	if err != nil {
		ctx.Logger().Error("Failed to get latest articles: ", err)
		errorMsg := presenter.ConvertErrorMessage(err)
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Error:  "Failed to get articles",
			Detail: &errorMsg,
		})
	}

	return ctx.JSON(http.StatusOK, openapi.ArticlesResponse{
		Articles: presenter.ConvertArticles(output.Articles),
	})
}

func (h *APIHandler) GetArticlesByCategory(ctx echo.Context, slug string, params openapi.GetArticlesByCategoryParams) error {
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Error: "Category slug is required",
		})
	}

	page := 1
	if params.Page != nil {
		page = *params.Page
	}

	limit := 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	input := usecase.GetArticlesByCategoryUsecaseInput{
		CategorySlug: slug,
		Page:         page,
		Limit:        limit,
	}

	output, err := h.getArticlesByCategoryUsecase.Exec(ctx.Request().Context(), input)
	if err != nil {
		ctx.Logger().Error("Failed to get articles by category: ", err)
		errorMsg := presenter.ConvertErrorMessage(err)
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Error:  "Failed to get articles",
			Detail: &errorMsg,
		})
	}

	return ctx.JSON(http.StatusOK, openapi.ArticlesResponse{
		Articles:   presenter.ConvertArticles(output.Articles),
		Pagination: presenter.ConvertPagination(output.Pagination),
	})
}