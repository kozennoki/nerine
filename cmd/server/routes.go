package main

import (
	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/kozennoki/nerine/internal/infrastructure/config"
	"github.com/kozennoki/nerine/internal/interfaces/middleware"
	"github.com/labstack/echo/v4"
)

func setupRoutes(e *echo.Echo, di *DIContainer, cfg *config.Config) {
	// API key authentication middleware for generated routes
	apiKeyMiddleware := middleware.APIKeyAuth(cfg.NerineAPIKey)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/health" {
				return next(c)
			}
			return apiKeyMiddleware(next)(c)
		}
	})

	// Register OpenAPI generated routes
	openapi.RegisterHandlers(e, di.APIHandler)
}
