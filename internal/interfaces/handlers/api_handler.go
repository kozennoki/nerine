package handlers

import (
	"github.com/kozennoki/nerine/internal/openapi"
	"github.com/kozennoki/nerine/internal/usecase"
)

type APIHandler struct {
	getArticlesUsecase           usecase.GetArticlesUsecase
	getArticleByIDUsecase        usecase.GetArticleByIDUsecase
	getPopularArticlesUsecase    usecase.GetPopularArticlesUsecase
	getLatestArticlesUsecase     usecase.GetLatestArticlesUsecase
	getArticlesByCategoryUsecase usecase.GetArticlesByCategoryUsecase
	getCategoriesUsecase         usecase.GetCategoriesUsecase
}

func NewAPIHandler(
	getArticlesUsecase usecase.GetArticlesUsecase,
	getArticleByIDUsecase usecase.GetArticleByIDUsecase,
	getPopularArticlesUsecase usecase.GetPopularArticlesUsecase,
	getLatestArticlesUsecase usecase.GetLatestArticlesUsecase,
	getArticlesByCategoryUsecase usecase.GetArticlesByCategoryUsecase,
	getCategoriesUsecase usecase.GetCategoriesUsecase,
) *APIHandler {
	return &APIHandler{
		getArticlesUsecase:           getArticlesUsecase,
		getArticleByIDUsecase:        getArticleByIDUsecase,
		getPopularArticlesUsecase:    getPopularArticlesUsecase,
		getLatestArticlesUsecase:     getLatestArticlesUsecase,
		getArticlesByCategoryUsecase: getArticlesByCategoryUsecase,
		getCategoriesUsecase:         getCategoriesUsecase,
	}
}

// Compile-time check to ensure APIHandler implements openapi.ServerInterface
var _ openapi.ServerInterface = (*APIHandler)(nil)