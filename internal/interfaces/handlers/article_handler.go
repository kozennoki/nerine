package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/kozennoki/nerine/internal/usecase"
)

type ArticleHandler struct {
	getArticlesUsecase usecase.GetArticlesUsecase
}

func NewArticleHandler(
	getArticlesUsecase usecase.GetArticlesUsecase,
) *ArticleHandler {
	return &ArticleHandler{
		getArticlesUsecase: getArticlesUsecase,
	}
}

func (h *ArticleHandler) GetArticles(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	input := usecase.GetArticlesUsecaseInput{
		Limit:  limit,
		Offset: offset,
	}

	output, err := h.getArticlesUsecase.Exec(c.Request().Context(), input)
	if err != nil {
		// ログに詳細なエラーを出力
		c.Logger().Error("Failed to get articles: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get articles",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"articles": output.Articles,
	})
}