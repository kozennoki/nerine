package repository

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
)

type ArticleRepository interface {
	GetArticles(ctx context.Context, limit, offset int) ([]*entity.Article, error)
	GetArticleByID(ctx context.Context, id string) (*entity.Article, error)
	GetArticlesByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]*entity.Article, error)
	GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error)
	GetLatestArticles(ctx context.Context, limit int) ([]*entity.Article, error)
	CountArticles(ctx context.Context) (int, error)
	CountArticlesByCategory(ctx context.Context, categorySlug string) (int, error)
}