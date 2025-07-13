package handlers_test

import (
	"github.com/kozennoki/nerine/internal/interfaces/handlers"
	"github.com/kozennoki/nerine/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

// TestAPIHandlerMocks holds all mocks for APIHandler testing
type TestAPIHandlerMocks struct {
	GetArticlesUsecase           *mocks.MockGetArticlesUsecase
	GetArticleByIDUsecase        *mocks.MockGetArticleByIDUsecase
	GetPopularArticlesUsecase    *mocks.MockGetPopularArticlesUsecase
	GetLatestArticlesUsecase     *mocks.MockGetLatestArticlesUsecase
	GetArticlesByCategoryUsecase *mocks.MockGetArticlesByCategoryUsecase
	GetCategoriesUsecase         *mocks.MockGetCategoriesUsecase
	GetZennArticlesUsecase       *mocks.MockGetZennArticlesUsecase
}

// CreateTestAPIHandler creates APIHandler with mocks for testing
func CreateTestAPIHandler(ctrl *gomock.Controller) (*handlers.APIHandler, *TestAPIHandlerMocks) {
	mocks := &TestAPIHandlerMocks{
		GetArticlesUsecase:           mocks.NewMockGetArticlesUsecase(ctrl),
		GetArticleByIDUsecase:        mocks.NewMockGetArticleByIDUsecase(ctrl),
		GetPopularArticlesUsecase:    mocks.NewMockGetPopularArticlesUsecase(ctrl),
		GetLatestArticlesUsecase:     mocks.NewMockGetLatestArticlesUsecase(ctrl),
		GetArticlesByCategoryUsecase: mocks.NewMockGetArticlesByCategoryUsecase(ctrl),
		GetCategoriesUsecase:         mocks.NewMockGetCategoriesUsecase(ctrl),
		GetZennArticlesUsecase:       mocks.NewMockGetZennArticlesUsecase(ctrl),
	}

	handler := handlers.NewAPIHandler(
		mocks.GetArticlesUsecase,
		mocks.GetArticleByIDUsecase,
		mocks.GetPopularArticlesUsecase,
		mocks.GetLatestArticlesUsecase,
		mocks.GetArticlesByCategoryUsecase,
		mocks.GetCategoriesUsecase,
		mocks.GetZennArticlesUsecase,
	)

	return handler, mocks
}

// IntPtr is a helper function for creating int pointers
func IntPtr(i int) *int {
	return &i
}
