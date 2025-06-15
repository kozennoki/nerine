package microcms

import (
	"context"
	"fmt"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
	"github.com/microcmsio/microcms-go-sdk"
)

type articleRepository struct {
	microCMS *microcms.Client
}

func NewArticleRepository(apiKey, serviceID string) repository.ArticleRepository {
	client := microcms.New(serviceID, apiKey)
	return &articleRepository{
		microCMS: client,
	}
}

type article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type articleListResponse struct {
	Contents   []article `json:"contents"`
	TotalCount int       `json:"totalCount"`
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
}

func (r *articleRepository) GetArticles(ctx context.Context, limit, offset int) ([]*entity.Article, error) {
	var res articleListResponse
	params := microcms.ListParams{
		Endpoint: "blog",
		Limit:    limit,
		Offset:   offset,
	}

	err := r.microCMS.List(params, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}

	articles := make([]*entity.Article, len(res.Contents))
	for i, item := range res.Contents {
		articles[i] = &entity.Article{
			ID:          item.ID,
			Title:       item.Title,
			Image:       item.Image,
			Category:    item.Category,
			Description: item.Description,
			Body:        item.Body,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return articles, nil
}

func (r *articleRepository) GetArticleByID(ctx context.Context, id string) (*entity.Article, error) {
	var res article

	params := microcms.GetParams{
		Endpoint:  "blog",
		ContentID: id,
	}
	err := r.microCMS.Get(params, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to get article by ID: %w", err)
	}

	return &entity.Article{
		ID:          res.ID,
		Title:       res.Title,
		Image:       res.Image,
		Category:    res.Category,
		Description: res.Description,
		Body:        res.Body,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (r *articleRepository) GetArticlesByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]*entity.Article, error) {
	var res articleListResponse
	params := microcms.ListParams{
		Endpoint: "blog",
		Limit:    limit,
		Offset:   offset,
		Filters:  fmt.Sprintf("category[equals]%s", categorySlug),
	}

	err := r.microCMS.List(params, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles by category: %w", err)
	}

	articles := make([]*entity.Article, len(res.Contents))
	for i, item := range res.Contents {
		articles[i] = &entity.Article{
			ID:          item.ID,
			Title:       item.Title,
			Image:       item.Image,
			Category:    item.Category,
			Description: item.Description,
			Body:        item.Body,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return articles, nil
}

func (r *articleRepository) GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error) {
	return r.GetArticles(ctx, limit, 0)
}

func (r *articleRepository) GetLatestArticles(ctx context.Context, limit int) ([]*entity.Article, error) {
	var res articleListResponse
	params := microcms.ListParams{
		Endpoint: "blog",
		Limit:    limit,
		Offset:   0,
		Orders:   []string{"-createdAt"},
	}

	err := r.microCMS.List(params, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest articles: %w", err)
	}

	articles := make([]*entity.Article, len(res.Contents))
	for i, item := range res.Contents {
		articles[i] = &entity.Article{
			ID:          item.ID,
			Title:       item.Title,
			Image:       item.Image,
			Category:    item.Category,
			Description: item.Description,
			Body:        item.Body,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return articles, nil
}

func (r *articleRepository) CountArticles(ctx context.Context) (int, error) {
	var res articleListResponse
	params := microcms.ListParams{
		Endpoint: "blog",
		Limit:    0,
	}

	err := r.microCMS.List(params, &res)
	if err != nil {
		return 0, fmt.Errorf("failed to count articles: %w", err)
	}

	return res.TotalCount, nil
}

func (r *articleRepository) CountArticlesByCategory(ctx context.Context, categorySlug string) (int, error) {
	var res articleListResponse
	params := microcms.ListParams{
		Endpoint: "blog",
		Limit:    0,
		Filters:  fmt.Sprintf("category[equals]%s", categorySlug),
	}

	err := r.microCMS.List(params, &res)
	if err != nil {
		return 0, fmt.Errorf("failed to count articles by category: %w", err)
	}

	return res.TotalCount, nil
}
