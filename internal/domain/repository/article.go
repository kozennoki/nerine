package repository

import (
	"context"

	"github.com/kozennoki/nerine/internal/domain/entity"
)

// ArticleReader は一覧取得の最小インターフェース。
// Zenn はこれだけ実装する。
type ArticleReader interface {
	GetArticles(ctx context.Context, limit, offset int) ([]*entity.Article, error)
}

// ArticleAdvancedReader は microCMS 側が提供する拡張機能。
type ArticleAdvancedReader interface {
	GetArticleByID(ctx context.Context, id string) (*entity.Article, error)
	GetArticlesByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]*entity.Article, error)
	GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error)
	GetLatestArticles(ctx context.Context, limit int) ([]*entity.Article, error)
	CountArticles(ctx context.Context) (int, error)
	CountArticlesByCategory(ctx context.Context, categorySlug string) (int, error)
}

type ArticleRepository interface {
	ArticleReader
	ArticleAdvancedReader
}
