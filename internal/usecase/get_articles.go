package usecase

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
	"github.com/kozennoki/nerine/internal/infrastructure/utils"
)

type GetArticlesUsecase interface {
	Exec(ctx context.Context, input GetArticlesUsecaseInput) (GetArticlesUsecaseOutput, error)
}

type GetArticlesUsecaseInput struct {
	Page  int
	Limit int
}

type GetArticlesUsecaseOutput struct {
	Articles   []*entity.Article
	Pagination utils.Pagination
}

type getArticles struct {
	articleRepo repository.ArticleRepository
}

func NewGetArticles(
	articleRepo repository.ArticleRepository,
) GetArticlesUsecase {
	return &getArticles{
		articleRepo: articleRepo,
	}
}

func (u *getArticles) Exec(
	ctx context.Context,
	input GetArticlesUsecaseInput,
) (GetArticlesUsecaseOutput, error) {
	// Get total count for pagination
	total, err := u.articleRepo.CountArticles(ctx)
	if err != nil {
		return GetArticlesUsecaseOutput{}, err
	}

	// Validate pagination parameters
	limit, offset, pagination := BuildPagination(
		input.Page, input.Limit, 10, 100, total,
	)

	// Get articles
	articles, err := u.articleRepo.GetArticles(ctx, limit, offset)
	if err != nil {
		return GetArticlesUsecaseOutput{}, err
	}

	return GetArticlesUsecaseOutput{
		Articles:   articles,
		Pagination: pagination,
	}, nil
}
