package microcms

import (
	"context"
	"fmt"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
	"github.com/microcmsio/microcms-go-sdk"
)

type categoryRepository struct {
	microCMS *microcms.Client
}

func NewCategoryRepository(apiKey, serviceID string) repository.CategoryRepository {
	client := microcms.New(serviceID, apiKey)
	return &categoryRepository{
		microCMS: client,
	}
}

type category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type categoryListResponse struct {
	Contents   []category `json:"contents"`
	TotalCount int        `json:"totalCount"`
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
}

func (r *categoryRepository) GetCategories(
	ctx context.Context,
) ([]*entity.Category, error) {
	var res categoryListResponse
	params := microcms.ListParams{
		Endpoint: "categories",
	}

	err := r.microCMS.List(params, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	categories := make([]*entity.Category, len(res.Contents))
	for i, item := range res.Contents {
		categories[i] = &entity.Category{
			Slug: item.ID,
			Name: item.Name,
		}
	}

	return categories, nil
}

func (r *categoryRepository) GetCategoryBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var res category

	params := microcms.GetParams{
		Endpoint:  "categories",
		ContentID: slug,
	}
	err := r.microCMS.Get(params, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	return &entity.Category{
		Slug: res.ID,
		Name: res.Name,
	}, nil
}
