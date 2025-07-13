package handlers

import (
	"net/http"

	"github.com/kozennoki/nerine/internal/interfaces/presenter"
	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/labstack/echo/v4"
)

func (h *APIHandler) GetZennArticles(c echo.Context, params openapi.GetZennArticlesParams) error {
	page := 1
	limit := 10

	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	input := usecase.GetZennArticlesUsecaseInput{
		Page:  page,
		Limit: limit,
	}

	output, err := h.getZennArticlesUsecase.Exec(c.Request().Context(), input)
	if err != nil {
		errMsg := presenter.ConvertErrorMessage(err)
		return c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Error:  "Failed to get Zenn articles",
			Detail: &errMsg,
		})
	}

	articles := presenter.ConvertArticles(output.Articles)
	pagination := presenter.ConvertPagination(output.Pagination)

	response := openapi.ArticlesResponse{
		Articles:   articles,
		Pagination: pagination,
	}

	return c.JSON(http.StatusOK, response)
}
