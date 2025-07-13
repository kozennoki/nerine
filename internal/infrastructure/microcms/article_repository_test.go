package microcms_test

import (
	"testing"

	"github.com/kozennoki/nerine/internal/infrastructure/microcms"
)

func TestNewArticleRepository(t *testing.T) {
	t.Parallel()

	apiKey := "test-api-key"
	serviceID := "test-service-id"

	repo := microcms.NewArticleRepository(apiKey, serviceID)

	if repo == nil {
		t.Error("NewArticleRepository() returned nil")
	}
}