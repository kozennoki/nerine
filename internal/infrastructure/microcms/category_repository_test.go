package microcms_test

import (
	"testing"

	"github.com/kozennoki/nerine/internal/infrastructure/microcms"
)

func TestNewCategoryRepository(t *testing.T) {
	t.Parallel()

	apiKey := "test-api-key"
	serviceID := "test-service-id"

	repo := microcms.NewCategoryRepository(apiKey, serviceID)

	if repo == nil {
		t.Error("NewCategoryRepository() returned nil")
	}
}
