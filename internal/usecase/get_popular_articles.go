package usecase

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
)

type GetPopularArticlesUsecase interface {
	Exec(context.Context, GetPopularArticlesUsecaseInput) (GetPopularArticlesUsecaseOutput, error)
}

type GetPopularArticlesUsecaseInput struct {
	Limit int
}

type GetPopularArticlesUsecaseOutput struct {
	Articles []*entity.Article
}

type getPopularArticles struct {
	repo repository.ArticleRepository
}

func NewGetPopularArticles(
	repo repository.ArticleRepository,
) GetPopularArticlesUsecase {
	return &getPopularArticles{
		repo: repo,
	}
}

func (u *getPopularArticles) Exec(
	ctx context.Context,
	input GetPopularArticlesUsecaseInput,
) (GetPopularArticlesUsecaseOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 5 // default limit for popular articles
	}
	if limit > 20 {
		limit = 20 // max limit as per OpenAPI spec
	}

	articles, err := u.repo.GetPopularArticles(ctx, limit)
	if err != nil {
		return GetPopularArticlesUsecaseOutput{}, err
	}

	return GetPopularArticlesUsecaseOutput{
		Articles: articles,
	}, nil
}
