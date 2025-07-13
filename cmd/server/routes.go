package main

import (
	"net/http"

	"github.com/kozennoki/nerine/internal/infrastructure/config"
	"github.com/kozennoki/nerine/internal/interfaces/middleware"
	"github.com/labstack/echo/v4"
)

func setupRoutes(e *echo.Echo, di *DIContainer, cfg *config.Config) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api := e.Group("/api/v1")
	api.Use(middleware.APIKeyAuth(cfg.NerineAPIKey))
	api.GET("/articles", di.ArticleHandler.GetArticles)
	api.GET("/articles/:id", di.ArticleHandler.GetArticleByID)
	api.GET("/articles/popular", di.ArticleHandler.GetPopularArticles)
	api.GET("/articles/latest", di.ArticleHandler.GetLatestArticles)
	api.GET("/categories", di.CategoryHandler.GetCategories)
	api.GET("/categories/:slug/articles", di.ArticleHandler.GetArticlesByCategory)
}
