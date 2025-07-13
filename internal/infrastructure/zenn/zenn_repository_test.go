package zenn_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kozennoki/nerine/internal/infrastructure/zenn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZennRepository_GetArticles_Success(t *testing.T) {
	t.Parallel()

	publishedAt, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2023-01-02T00:00:00Z")

	mockResponse := map[string]interface{}{
		"articles": []map[string]interface{}{
			{
				"id":                123,
				"post_type":         "Article",
				"title":             "Test Article",
				"slug":              "test-article",
				"comments_count":    5,
				"liked_count":       10,
				"bookmarked_count":  3,
				"body_letters_count": 1000,
				"article_type":      "tech",
				"emoji":             "📝",
				"published_at":      publishedAt.Format(time.RFC3339),
				"body_updated_at":   updatedAt.Format(time.RFC3339),
				"user": map[string]interface{}{
					"id":               1,
					"username":         "kozennoki",
					"name":             "Test User",
					"avatar_small_url": "https://example.com/avatar.jpg",
				},
			},
		},
		"next_page":   nil,
		"total_count": nil,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/articles", r.URL.Path)
		assert.Equal(t, "kozennoki", r.URL.Query().Get("username"))
		assert.Equal(t, "latest", r.URL.Query().Get("order"))
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	repo := zenn.NewZennRepositoryWithBaseURL(server.URL)
	
	articles, err := repo.GetArticles(context.Background(), 10, 0)
	
	require.NoError(t, err)
	require.Len(t, articles, 1)
	
	article := articles[0]
	assert.Equal(t, "123", article.ID)
	assert.Equal(t, "Test Article", article.Title)
	assert.Equal(t, "📝", article.Image)
	assert.Equal(t, "zenn", article.Category.Slug)
	assert.Equal(t, "Zenn", article.Category.Name)
	assert.Equal(t, "Zenn記事 - test-article", article.Description)
	assert.Equal(t, "", article.Body)
	assert.Equal(t, publishedAt.UTC(), article.CreatedAt)
	assert.Equal(t, updatedAt.UTC(), article.UpdatedAt)
}

func TestZennRepository_GetArticles_PaginationCalculation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		limit          int
		offset         int
		expectedPage   string
	}{
		{
			name:         "First page",
			limit:        10,
			offset:       0,
			expectedPage: "1",
		},
		{
			name:         "Second page",
			limit:        10,
			offset:       10,
			expectedPage: "2",
		},
		{
			name:         "Third page with different limit",
			limit:        5,
			offset:       10,
			expectedPage: "3",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockResponse := map[string]interface{}{
				"articles":    []map[string]interface{}{},
				"next_page":   nil,
				"total_count": nil,
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPage, r.URL.Query().Get("page"))
				assert.Equal(t, "kozennoki", r.URL.Query().Get("username"))
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(mockResponse)
			}))
			defer server.Close()

			repo := zenn.NewZennRepositoryWithBaseURL(server.URL)
			_, err := repo.GetArticles(context.Background(), tc.limit, tc.offset)
			
			assert.NoError(t, err)
		})
	}
}

func TestZennRepository_GetArticles_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	repo := zenn.NewZennRepositoryWithBaseURL(server.URL)
	
	articles, err := repo.GetArticles(context.Background(), 10, 0)
	
	assert.Error(t, err)
	assert.Nil(t, articles)
	assert.Contains(t, err.Error(), "zenn API returned status 500")
}

func TestZennRepository_GetArticles_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	repo := zenn.NewZennRepositoryWithBaseURL(server.URL)
	
	articles, err := repo.GetArticles(context.Background(), 10, 0)
	
	assert.Error(t, err)
	assert.Nil(t, articles)
	assert.Contains(t, err.Error(), "failed to decode Zenn response")
}

func TestZennRepository_GetArticles_ContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	repo := zenn.NewZennRepositoryWithBaseURL(server.URL)
	
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	articles, err := repo.GetArticles(ctx, 10, 0)
	
	assert.Error(t, err)
	assert.Nil(t, articles)
	assert.Contains(t, err.Error(), "failed to fetch articles from Zenn")
}

func TestZennRepository_GetArticles_ZeroLimit(t *testing.T) {
	t.Parallel()

	repo := zenn.NewZennRepositoryWithBaseURL("http://example.com")
	
	// This should handle the zero limit case
	articles, err := repo.GetArticles(context.Background(), 0, 0)
	
	assert.Error(t, err)
	assert.Nil(t, articles)
	assert.Contains(t, err.Error(), "limit must be greater than 0")
}

func TestZennRepository_GetArticles_InvalidURL(t *testing.T) {
	t.Parallel()

	// Use an invalid URL that would cause http.NewRequestWithContext to fail
	repo := zenn.NewZennRepositoryWithBaseURL("ht\ttp://invalid-url")
	
	articles, err := repo.GetArticles(context.Background(), 10, 0)
	
	assert.Error(t, err)
	assert.Nil(t, articles)
	assert.Contains(t, err.Error(), "failed to create request")
}

func TestZennRepository_NotImplementedMethods(t *testing.T) {
	t.Parallel()

	repo := zenn.NewZennRepository()
	ctx := context.Background()

	t.Run("GetArticleByID", func(t *testing.T) {
		article, err := repo.GetArticleByID(ctx, "123")
		assert.Error(t, err)
		assert.Nil(t, article)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("GetArticlesByCategory", func(t *testing.T) {
		articles, err := repo.GetArticlesByCategory(ctx, "tech", 10, 0)
		assert.Error(t, err)
		assert.Nil(t, articles)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("GetPopularArticles", func(t *testing.T) {
		articles, err := repo.GetPopularArticles(ctx, 10)
		assert.Error(t, err)
		assert.Nil(t, articles)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("GetLatestArticles", func(t *testing.T) {
		articles, err := repo.GetLatestArticles(ctx, 10)
		assert.Error(t, err)
		assert.Nil(t, articles)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("CountArticles", func(t *testing.T) {
		count, err := repo.CountArticles(ctx)
		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("CountArticlesByCategory", func(t *testing.T) {
		count, err := repo.CountArticlesByCategory(ctx, "tech")
		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Contains(t, err.Error(), "not implemented")
	})
}