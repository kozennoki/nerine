package presenter_test

import (
	"errors"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/infrastructure/utils"
	"github.com/kozennoki/nerine/internal/interfaces/presenter"
	"github.com/kozennoki/nerine/internal/openapi"
)

func TestConvertArticle(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	article := &entity.Article{
		ID:    "test-id",
		Title: "Test Title",
		Image: "https://example.com/image.jpg",
		Category: entity.Category{
			Slug: "tech",
			Name: "Technology",
		},
		Description: "Test Description",
		Body:        "Test Body",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	result := presenter.ConvertArticle(article)

	expected := openapi.Article{
		ID:    "test-id",
		Title: "Test Title",
		Image: "https://example.com/image.jpg",
		Category: openapi.Category{
			Slug: "tech",
			Name: "Technology",
		},
		Description: "Test Description",
		Body:        "Test Body",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	if result != expected {
		t.Errorf("ConvertArticle() = %+v, want %+v", result, expected)
	}
}

func TestConvertArticles(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	articles := []*entity.Article{
		{
			ID:    "test-id-1",
			Title: "Test Title 1",
			Image: "https://example.com/image1.jpg",
			Category: entity.Category{
				Slug: "tech",
				Name: "Technology",
			},
			Description: "Test Description 1",
			Body:        "Test Body 1",
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
		{
			ID:    "test-id-2",
			Title: "Test Title 2",
			Image: "https://example.com/image2.jpg",
			Category: entity.Category{
				Slug: "design",
				Name: "Design",
			},
			Description: "Test Description 2",
			Body:        "Test Body 2",
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
	}

	result := presenter.ConvertArticles(articles)

	if len(result) != 2 {
		t.Errorf("ConvertArticles() returned %d articles, want 2", len(result))
	}

	if result[0].ID != "test-id-1" {
		t.Errorf("ConvertArticles()[0].ID = %s, want test-id-1", result[0].ID)
	}

	if result[1].ID != "test-id-2" {
		t.Errorf("ConvertArticles()[1].ID = %s, want test-id-2", result[1].ID)
	}
}

func TestConvertArticles_EmptySlice(t *testing.T) {
	t.Parallel()

	articles := []*entity.Article{}
	result := presenter.ConvertArticles(articles)

	if len(result) != 0 {
		t.Errorf("ConvertArticles() with empty slice returned %d articles, want 0", len(result))
	}
}

func TestConvertCategory(t *testing.T) {
	t.Parallel()

	category := entity.Category{
		Slug: "tech",
		Name: "Technology",
	}

	result := presenter.ConvertCategory(category)

	expected := openapi.Category{
		Slug: "tech",
		Name: "Technology",
	}

	if result != expected {
		t.Errorf("ConvertCategory() = %+v, want %+v", result, expected)
	}
}

func TestConvertCategories(t *testing.T) {
	t.Parallel()

	categories := []*entity.Category{
		{
			Slug: "tech",
			Name: "Technology",
		},
		{
			Slug: "design",
			Name: "Design",
		},
	}

	result := presenter.ConvertCategories(categories)

	if len(result) != 2 {
		t.Errorf("ConvertCategories() returned %d categories, want 2", len(result))
	}

	if result[0].Slug != "tech" {
		t.Errorf("ConvertCategories()[0].Slug = %s, want tech", result[0].Slug)
	}

	if result[1].Slug != "design" {
		t.Errorf("ConvertCategories()[1].Slug = %s, want design", result[1].Slug)
	}
}

func TestConvertCategories_EmptySlice(t *testing.T) {
	t.Parallel()

	categories := []*entity.Category{}
	result := presenter.ConvertCategories(categories)

	if len(result) != 0 {
		t.Errorf("ConvertCategories() with empty slice returned %d categories, want 0", len(result))
	}
}

func TestConvertPagination(t *testing.T) {
	t.Parallel()

	pagination := utils.Pagination{
		Total:      100,
		Page:       2,
		Limit:      10,
		TotalPages: 10,
	}

	result := presenter.ConvertPagination(pagination)

	if result.Total == nil || *result.Total != 100 {
		t.Errorf("ConvertPagination().Total = %v, want 100", result.Total)
	}

	if result.Page == nil || *result.Page != 2 {
		t.Errorf("ConvertPagination().Page = %v, want 2", result.Page)
	}

	if result.Limit == nil || *result.Limit != 10 {
		t.Errorf("ConvertPagination().Limit = %v, want 10", result.Limit)
	}

	if result.TotalPages == nil || *result.TotalPages != 10 {
		t.Errorf("ConvertPagination().TotalPages = %v, want 10", result.TotalPages)
	}
}

func TestConvertErrorMessage(t *testing.T) {
	t.Parallel()

	err := errors.New("test error message")
	result := presenter.ConvertErrorMessage(err)

	expected := "test error message"
	if result != expected {
		t.Errorf("ConvertErrorMessage() = %s, want %s", result, expected)
	}
}
