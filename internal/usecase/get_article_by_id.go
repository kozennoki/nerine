package usecase

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
)

type GetArticleByIDUsecase interface {
	Exec(ctx context.Context, input GetArticleByIDUsecaseInput) (GetArticleByIDUsecaseOutput, error)
}

type GetArticleByIDUsecaseInput struct {
	ID string
}

type GetArticleByIDUsecaseOutput struct {
	Article *entity.Article
}

type getArticleByID struct {
	articleRepo repository.ArticleRepository
}

func NewGetArticleByID(
	articleRepo repository.ArticleRepository,
) GetArticleByIDUsecase {
	return &getArticleByID{
		articleRepo: articleRepo,
	}
}

func (u *getArticleByID) Exec(ctx context.Context, input GetArticleByIDUsecaseInput) (GetArticleByIDUsecaseOutput, error) {
	article, err := u.articleRepo.GetArticleByID(ctx, input.ID)
	if err != nil {
		return GetArticleByIDUsecaseOutput{}, err
	}

	return GetArticleByIDUsecaseOutput{
		Article: article,
	}, nil
}