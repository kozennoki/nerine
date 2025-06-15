package repository

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]*entity.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*entity.Category, error)
}