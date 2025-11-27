package zenn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository"
)

const (
	baseURL  = "https://zenn.dev/api"
	userName = "kozennoki"
)

type zennRepository struct {
	httpClient *http.Client
	baseURL    string
}

func NewZennRepository() repository.ArticleReader {
	return &zennRepository{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func NewZennRepositoryWithBaseURL(baseURL string) repository.ArticleReader {
	return &zennRepository{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

type zennArticle struct {
	ID               int       `json:"id"`
	PostType         string    `json:"post_type"`
	Title            string    `json:"title"`
	Slug             string    `json:"slug"`
	CommentsCount    int       `json:"comments_count"`
	LikedCount       int       `json:"liked_count"`
	BookmarkedCount  int       `json:"bookmarked_count"`
	BodyLettersCount int       `json:"body_letters_count"`
	ArticleType      string    `json:"article_type"`
	Emoji            string    `json:"emoji"`
	PublishedAt      time.Time `json:"published_at"`
	UpdatedAt        time.Time `json:"body_updated_at"`
	User             zennUser  `json:"user"`
}

type zennUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	AvatarS3 string `json:"avatar_small_url"`
}

type zennAPIResponse struct {
	Articles   []zennArticle `json:"articles"`
	NextPage   *int          `json:"next_page"`
	TotalCount *int          `json:"total_count"`
}

func (r *zennRepository) GetArticles(ctx context.Context, limit, offset int) ([]*entity.Article, error) {
	if limit == 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}
	page := (offset / limit) + 1

	url := fmt.Sprintf(r.baseURL+"/articles?username=%s&order=latest&page=%d", userName, page)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles from Zenn: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zenn API returned status %d", resp.StatusCode)
	}

	var zennResp zennAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&zennResp); err != nil {
		return nil, fmt.Errorf("failed to decode Zenn response: %w", err)
	}

	articles := make([]*entity.Article, 0, len(zennResp.Articles))
	for _, zennArticle := range zennResp.Articles {
		article := convertToEntity(zennArticle)
		articles = append(articles, article)
	}

	return articles, nil
}

func convertToEntity(zennArticle zennArticle) *entity.Article {
	return &entity.Article{
		ID:    zennArticle.Slug,
		Title: zennArticle.Emoji + zennArticle.Title,
		Category: entity.Category{
			Slug: "zenn",
			Name: "Zenn",
		},
		Description: fmt.Sprintf("Zenn記事 - %s", strconv.Itoa(zennArticle.ID)),
		Body:        "",
		PublishedAt: zennArticle.PublishedAt.UTC(),
		CreatedAt:   zennArticle.PublishedAt.UTC(),
		UpdatedAt:   zennArticle.UpdatedAt.UTC(),
	}
}
