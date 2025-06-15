package main

import (
	"log"
	"net/http"

	"github.com/kozennoki/nerine/internal/infrastructure/config"
	"github.com/kozennoki/nerine/internal/infrastructure/logger"
	"github.com/kozennoki/nerine/internal/infrastructure/microcms"
	"github.com/kozennoki/nerine/internal/interfaces/handlers"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := logger.New()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()

	cfg, err := config.Load()
	if err != nil {
		zapLogger.Fatal("Failed to load config", zap.Error(err))
	}

	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Repository
	articleRepo := microcms.NewArticleRepository(cfg.MicroCMSAPIKey, cfg.MicroCMSServiceID)
	categoryRepo := microcms.NewCategoryRepository(cfg.MicroCMSAPIKey, cfg.MicroCMSServiceID)

	// UseCase
	getArticlesUsecase := usecase.NewGetArticles(articleRepo)
	getArticleByIDUsecase := usecase.NewGetArticleByID(articleRepo)
	getCategoriesUsecase := usecase.NewGetCategories(categoryRepo)

	// Handler
	articleHandler := handlers.NewArticleHandler(getArticlesUsecase, getArticleByIDUsecase)
	categoryHandler := handlers.NewCategoryHandler(getCategoriesUsecase)

	// Routes
	api := e.Group("/api/v1")
	// api.Use(middleware.APIKeyAuth(cfg.NerineAPIKey))
	api.GET("/articles", articleHandler.GetArticles)
	api.GET("/articles/:id", articleHandler.GetArticleByID)
	api.GET("/categories", categoryHandler.GetCategories)

	zapLogger.Info("Starting server", zap.String("port", cfg.Port))
	if err := e.Start(":" + cfg.Port); err != nil {
		zapLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
