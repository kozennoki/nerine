package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func APIKeyAuth(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestAPIKey := c.Request().Header.Get("X-API-Key")
			if requestAPIKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing API key")
			}

			if requestAPIKey != apiKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
			}

			return next(c)
		}
	}
}