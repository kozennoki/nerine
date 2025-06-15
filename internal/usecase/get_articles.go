package usecase

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
)

type GetArticlesUsecase interface {
	Exec(ctx context.Context, input GetArticlesUsecaseInput) (GetArticlesUsecaseOutput, error)
}

type GetArticlesUsecaseInput struct {
	Limit  int
	Offset int
}

type GetArticlesUsecaseOutput struct {
	Articles []*entity.Article
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

func (u *getArticles) Exec(ctx context.Context, input GetArticlesUsecaseInput) (GetArticlesUsecaseOutput, error) {
	articles, err := u.articleRepo.GetArticles(ctx, input.Limit, input.Offset)
	if err != nil {
		return GetArticlesUsecaseOutput{}, err
	}

	return GetArticlesUsecaseOutput{
		Articles: articles,
	}, nil
}