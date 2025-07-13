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
	limit := ValidateLimit(input.Limit, 5, 20)

	articles, err := u.repo.GetLatestArticles(ctx, limit)
	if err != nil {
		return GetLatestArticlesUsecaseOutput{}, err
	}

	return GetLatestArticlesUsecaseOutput{
		Articles: articles,
	}, nil
}
