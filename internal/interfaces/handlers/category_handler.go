package handlers

import (
	"net/http"

	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/kozennoki/nerine/internal/interfaces/presenter"
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

func (h *CategoryHandler) GetCategories(ctx echo.Context) error {
	input := usecase.GetCategoriesUsecaseInput{}

	output, err := h.getCategoriesUsecase.Exec(ctx.Request().Context(), input)
	if err != nil {
		ctx.Logger().Error("Failed to get categories: ", err)
		errorMsg := presenter.ConvertErrorMessage(err)
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Error:  "Failed to get categories",
			Detail: &errorMsg,
		})
	}

	return ctx.JSON(http.StatusOK, openapi.CategoriesResponse{
		Categories: presenter.ConvertCategories(output.Categories),
	})
}