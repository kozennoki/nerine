package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/kozennoki/nerine/internal/usecase"
)

type CategoryHandler struct {
	getCategoriesUsecase usecase.GetCategoriesUsecase
}

func NewCategoryHandler(
	getCategoriesUsecase usecase.GetCategoriesUsecase,
) *CategoryHandler {
	return &CategoryHandler{
		getCategoriesUsecase: getCategoriesUsecase,
	}
}

func (h *CategoryHandler) GetCategories(c echo.Context) error {
	input := usecase.GetCategoriesUsecaseInput{}

	output, err := h.getCategoriesUsecase.Exec(c.Request().Context(), input)
	if err != nil {
		c.Logger().Error("Failed to get categories: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":  "Failed to get categories",
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": output.Categories,
	})
}