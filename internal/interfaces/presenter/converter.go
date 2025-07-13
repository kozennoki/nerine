package presenter

import (
	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/infrastructure/utils"
	"github.com/kozennoki/nerine/internal/openapi"
)

func ConvertArticle(article *entity.Article) openapi.Article {
	return openapi.Article{
		ID:          article.ID,
		Title:       article.Title,
		Image:       article.Image,
		Category:    ConvertCategory(article.Category),
		Description: article.Description,
		Body:        article.Body,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
}

func ConvertArticles(articles []*entity.Article) []openapi.Article {
	result := make([]openapi.Article, len(articles))
	for i, article := range articles {
		result[i] = ConvertArticle(article)
	}
	return result
}

func ConvertCategory(category entity.Category) openapi.Category {
	return openapi.Category{
		Slug: category.Slug,
		Name: category.Name,
	}
}

func ConvertCategories(categories []*entity.Category) []openapi.Category {
	result := make([]openapi.Category, len(categories))
	for i, category := range categories {
		result[i] = ConvertCategory(*category)
	}
	return result
}

func ConvertPagination(pagination utils.Pagination) *openapi.Pagination {
	total := pagination.Total
	page := pagination.Page
	limit := pagination.Limit
	totalPages := pagination.TotalPages

	return &openapi.Pagination{
		Total:      &total,
		Page:       &page,
		Limit:      &limit,
		TotalPages: &totalPages,
	}
}

func ConvertErrorMessage(err error) string {
	return err.Error()
}
