package usecase

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
)

type GetLatestArticlesUsecase interface {
	Exec(context.Context, GetLatestArticlesUsecaseInput) (GetLatestArticlesUsecaseOutput, error)
}

type GetLatestArticlesUsecaseInput struct {
	Limit int
}

type GetLatestArticlesUsecaseOutput struct {
	Articles []*entity.Article
}

type getLatestArticles struct {
	repo repository.ArticleRepository
}

func NewGetLatestArticles(
	repo repository.ArticleRepository,
) GetLatestArticlesUsecase {
	return &getLatestArticles{
		repo: repo,
	}
}

func (u *getLatestArticles) Exec(
	ctx context.Context,
	input GetLatestArticlesUsecaseInput,
) (GetLatestArticlesUsecaseOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 5 // default limit for latest articles
	}
	if limit > 20 {
		limit = 20 // max limit as per OpenAPI spec
	}

	articles, err := u.repo.GetLatestArticles(ctx, limit)
	if err != nil {
		return GetLatestArticlesUsecaseOutput{}, err
	}

	return GetLatestArticlesUsecaseOutput{
		Articles: articles,
	}, nil
}
