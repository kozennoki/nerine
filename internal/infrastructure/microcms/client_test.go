package microcms_test

import (
	"testing"

	"github.com/kozennoki/nerine/internal/infrastructure/microcms"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	apiKey := "test-api-key"
	serviceID := "test-service-id"

	client := microcms.NewClient(apiKey, serviceID)

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
}

func TestNewClient_WithEmptyValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		apiKey    string
		serviceID string
	}{
		{
			name:      "Empty API key",
			apiKey:    "",
			serviceID: "test-service-id",
		},
		{
			name:      "Empty service ID",
			apiKey:    "test-api-key",
			serviceID: "",
		},
		{
			name:      "Both empty",
			apiKey:    "",
			serviceID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := microcms.NewClient(tt.apiKey, tt.serviceID)

			if client == nil {
				t.Error("Expected client to be non-nil even with empty values")
			}
		})
	}
}
