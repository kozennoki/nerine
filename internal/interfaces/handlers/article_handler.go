package handlers

import (
	"net/http"
	"strconv"

	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	getArticlesUsecase           usecase.GetArticlesUsecase
	getArticleByIDUsecase        usecase.GetArticleByIDUsecase
	getPopularArticlesUsecase    usecase.GetPopularArticlesUsecase
	getLatestArticlesUsecase     usecase.GetLatestArticlesUsecase
	getArticlesByCategoryUsecase usecase.GetArticlesByCategoryUsecase
}

func NewArticleHandler(
	getArticlesUsecase usecase.GetArticlesUsecase,
	getArticleByIDUsecase usecase.GetArticleByIDUsecase,
	getPopularArticlesUsecase usecase.GetPopularArticlesUsecase,
	getLatestArticlesUsecase usecase.GetLatestArticlesUsecase,
	getArticlesByCategoryUsecase usecase.GetArticlesByCategoryUsecase,
) *ArticleHandler {
	return &ArticleHandler{
		getArticlesUsecase:           getArticlesUsecase,
		getArticleByIDUsecase:        getArticleByIDUsecase,
		getPopularArticlesUsecase:    getPopularArticlesUsecase,
		getLatestArticlesUsecase:     getLatestArticlesUsecase,
		getArticlesByCategoryUsecase: getArticlesByCategoryUsecase,
	}
}

func (h *ArticleHandler) GetArticles(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 1 {
			page = p
		}
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	input := usecase.GetArticlesUsecaseInput{
		Page:  page,
		Limit: limit,
	}

	output, err := h.getArticlesUsecase.Exec(c.Request().Context(), input)
	if err != nil {
		c.Logger().Error("Failed to get articles: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":  "Failed to get articles",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"articles":   output.Articles,
		"pagination": output.Pagination,
	})
}

func (h *ArticleHandler) GetArticleByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Article ID is required",
		})
	}

	input := usecase.GetArticleByIDUsecaseInput{
		ID: id,
	}

	output, err := h.getArticleByIDUsecase.Exec(c.Request().Context(), input)
	if err != nil {
		c.Logger().Error("Failed to get article by ID: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":  "Failed to get article",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"article": output.Article,
	})
}

func (h *ArticleHandler) GetPopularArticles(c echo.Context) error {
	limitStr := c.QueryParam("limit")

	limit := 5
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	input := usecase.GetPopularArticlesUsecaseInput{
		Limit: limit,
	}

	output, err := h.getPopularArticlesUsecase.Exec(c.Request().Context(), input)
	if err != nil {
		c.Logger().Error("Failed to get popular articles: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":  "Failed to get articles",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"articles": output.Articles,
	})
}

func (h *ArticleHandler) GetLatestArticles(c echo.Context) error {
	limitStr := c.QueryParam("limit")

	limit := 5
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	input := usecase.GetLatestArticlesUsecaseInput{
		Limit: limit,
	}

	output, err := h.getLatestArticlesUsecase.Exec(c.Request().Context(), input)
	if err != nil {
		c.Logger().Error("Failed to get latest articles: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":  "Failed to get articles",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"articles": output.Articles,
	})
}

func (h *ArticleHandler) GetArticlesByCategory(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Category slug is required",
		})
	}

	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 1 {
			page = p
		}
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	input := usecase.GetArticlesByCategoryUsecaseInput{
		CategorySlug: slug,
		Page:         page,
		Limit:        limit,
	}

	output, err := h.getArticlesByCategoryUsecase.Exec(c.Request().Context(), input)
	if err != nil {
		c.Logger().Error("Failed to get articles by category: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":  "Failed to get articles",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"articles":   output.Articles,
		"pagination": output.Pagination,
	})
}
