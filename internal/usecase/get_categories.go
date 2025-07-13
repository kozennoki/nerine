package usecase

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
)

type GetCategoriesUsecase interface {
	Exec(ctx context.Context, input GetCategoriesUsecaseInput) (GetCategoriesUsecaseOutput, error)
}

type GetCategoriesUsecaseInput struct{}

type GetCategoriesUsecaseOutput struct {
	Categories []*entity.Category
}

type getCategories struct {
	categoryRepo repository.CategoryRepository
}

func NewGetCategories(
	categoryRepo repository.CategoryRepository,
) GetCategoriesUsecase {
	return &getCategories{
		categoryRepo: categoryRepo,
	}
}

func (u *getCategories) Exec(
	ctx context.Context,
	input GetCategoriesUsecaseInput,
) (GetCategoriesUsecaseOutput, error) {
	categories, err := u.categoryRepo.GetCategories(ctx)
	if err != nil {
		return GetCategoriesUsecaseOutput{}, err
	}

	return GetCategoriesUsecaseOutput{
		Categories: categories,
	}, nil
}
