package handlers

import (
	"net/http"

	"github.com/kozennoki/nerine/internal/interfaces/presenter"
	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/labstack/echo/v4"
)

func (h *APIHandler) GetCategories(ctx echo.Context) error {
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
