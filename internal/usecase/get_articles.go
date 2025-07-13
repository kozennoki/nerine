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
	total, err := u.articleRepo.CountArticles(ctx)
	if err != nil {
		return GetArticlesUsecaseOutput{}, err
	}

	// Get articles
	articles, err := u.articleRepo.GetArticles(ctx, limit, offset)
	if err != nil {
		return GetArticlesUsecaseOutput{}, err
	}

	pagination := utils.NewPagination(total, page, limit)

	return GetArticlesUsecaseOutput{
		Articles:   articles,
		Pagination: pagination,
	}, nil
}
