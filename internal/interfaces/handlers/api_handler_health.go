package handlers

import (
	"net/http"

	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/labstack/echo/v4"
)

func (h *APIHandler) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, openapi.HealthResponse{
		Status: "ok",
	})
}