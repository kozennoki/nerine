package usecase

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
	"github.com/kozennoki/nerine/internal/infrastructure/utils"
)

type GetArticlesByCategoryUsecase interface {
	Exec(context.Context, GetArticlesByCategoryUsecaseInput) (GetArticlesByCategoryUsecaseOutput, error)
}

type GetArticlesByCategoryUsecaseInput struct {
	CategorySlug string
	Page         int
	Limit        int
}

type GetArticlesByCategoryUsecaseOutput struct {
	Articles   []*entity.Article
	Pagination utils.Pagination
}

type getArticlesByCategory struct {
	repo repository.ArticleRepository
}

func NewGetArticlesByCategory(
	repo repository.ArticleRepository,
) GetArticlesByCategoryUsecase {
	return &getArticlesByCategory{
		repo: repo,
	}
}

func (u *getArticlesByCategory) Exec(
	ctx context.Context,
	input GetArticlesByCategoryUsecaseInput,
) (GetArticlesByCategoryUsecaseOutput, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 10 // default limit
	}
	if limit > 100 {
		limit = 100 // max limit as per OpenAPI spec
	}

	offset := utils.ConvertPageToOffset(page, limit)

	// Get total count for pagination
	total, err := u.repo.CountArticlesByCategory(ctx, input.CategorySlug)
	if err != nil {
		return GetArticlesByCategoryUsecaseOutput{}, err
	}

	// Get articles
	articles, err := u.repo.GetArticlesByCategory(ctx, input.CategorySlug, limit, offset)
	if err != nil {
		return GetArticlesByCategoryUsecaseOutput{}, err
	}

	pagination := utils.NewPagination(total, page, limit)

	return GetArticlesByCategoryUsecaseOutput{
		Articles:   articles,
		Pagination: pagination,
	}, nil
}
