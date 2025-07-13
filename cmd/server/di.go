package main

import (
	"github.com/kozennoki/nerine/internal/infrastructure/config"
	"github.com/kozennoki/nerine/internal/infrastructure/microcms"
	"github.com/kozennoki/nerine/internal/interfaces/handlers"
	"github.com/kozennoki/nerine/internal/usecase"
)

type DIContainer struct {
	APIHandler *handlers.APIHandler
}

func NewDIContainer(cfg *config.Config) *DIContainer {
	// Repository
	articleRepo := microcms.NewArticleRepository(cfg.MicroCMSAPIKey, cfg.MicroCMSServiceID)
	categoryRepo := microcms.NewCategoryRepository(cfg.MicroCMSAPIKey, cfg.MicroCMSServiceID)

	// UseCase
	getArticlesUsecase := usecase.NewGetArticles(articleRepo)
	getArticleByIDUsecase := usecase.NewGetArticleByID(articleRepo)
	getPopularArticlesUsecase := usecase.NewGetPopularArticles(articleRepo)
	getLatestArticlesUsecase := usecase.NewGetLatestArticles(articleRepo)
	getArticlesByCategoryUsecase := usecase.NewGetArticlesByCategory(articleRepo)
	getCategoriesUsecase := usecase.NewGetCategories(categoryRepo)

	// Handler
	apiHandler := handlers.NewAPIHandler(
		getArticlesUsecase,
		getArticleByIDUsecase,
		getPopularArticlesUsecase,
		getLatestArticlesUsecase,
		getArticlesByCategoryUsecase,
		getCategoriesUsecase,
	)

	return &DIContainer{
		APIHandler: apiHandler,
	}
}