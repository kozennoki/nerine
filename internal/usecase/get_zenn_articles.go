package usecase

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
	"github.com/kozennoki/nerine/internal/infrastructure/utils"
)

type GetZennArticlesUsecase interface {
	Exec(ctx context.Context, input GetZennArticlesUsecaseInput) (GetZennArticlesUsecaseOutput, error)
}

type GetZennArticlesUsecaseInput struct {
	Page  int
	Limit int
}

type GetZennArticlesUsecaseOutput struct {
	Articles   []*entity.Article
	Pagination utils.Pagination
}

type getZennArticles struct {
	zennRepo repository.ArticleRepository
}

func NewGetZennArticles(
	zennRepo repository.ArticleRepository,
) GetZennArticlesUsecase {
	return &getZennArticles{
		zennRepo: zennRepo,
	}
}

func (u *getZennArticles) Exec(
	ctx context.Context,
	input GetZennArticlesUsecaseInput,
) (GetZennArticlesUsecaseOutput, error) {
	limit, offset, pagination := BuildPagination(
		input.Page, input.Limit, 10, 100, 0,
	)

	articles, err := u.zennRepo.GetArticles(ctx, limit, offset)
	if err != nil {
		return GetZennArticlesUsecaseOutput{}, err
	}

	return GetZennArticlesUsecaseOutput{
		Articles:   articles,
		Pagination: pagination,
	}, nil
}