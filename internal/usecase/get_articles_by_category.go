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
	// Get total count for pagination
	total, err := u.repo.CountArticlesByCategory(ctx, input.CategorySlug)
	if err != nil {
		return GetArticlesByCategoryUsecaseOutput{}, err
	}

	// Validate pagination parameters
	limit, offset, pagination := BuildPagination(
		input.Page, input.Limit, 10, 100, total,
	)

	// Get articles
	articles, err := u.repo.GetArticlesByCategory(ctx, input.CategorySlug, limit, offset)
	if err != nil {
		return GetArticlesByCategoryUsecaseOutput{}, err
	}

	return GetArticlesByCategoryUsecaseOutput{
		Articles:   articles,
		Pagination: pagination,
	}, nil
}
