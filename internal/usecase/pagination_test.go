package usecase_test

import (
	"testing"

	"github.com/kozennoki/nerine/internal/usecase"
)

func TestBuildPagination(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		page               int
		limit              int
		defaultLimit       int
		maxLimit           int
		total              int
		expectedPage       int
		expectedLimit      int
		expectedOffset     int
		expectedTotalPages int
	}{
		{
			name:               "Valid parameters",
			page:               2,
			limit:              10,
			defaultLimit:       10,
			maxLimit:           100,
			total:              50,
			expectedPage:       2,
			expectedLimit:      10,
			expectedOffset:     10,
			expectedTotalPages: 5,
		},
		{
			name:               "Page less than 1",
			page:               0,
			limit:              10,
			defaultLimit:       10,
			maxLimit:           100,
			total:              50,
			expectedPage:       1,
			expectedLimit:      10,
			expectedOffset:     0,
			expectedTotalPages: 5,
		},
		{
			name:               "Limit less than or equal to 0",
			page:               1,
			limit:              0,
			defaultLimit:       10,
			maxLimit:           100,
			total:              50,
			expectedPage:       1,
			expectedLimit:      10,
			expectedOffset:     0,
			expectedTotalPages: 5,
		},
		{
			name:               "Limit exceeds maximum",
			page:               1,
			limit:              150,
			defaultLimit:       10,
			maxLimit:           100,
			total:              50,
			expectedPage:       1,
			expectedLimit:      100,
			expectedOffset:     0,
			expectedTotalPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			limit, offset, pagination := usecase.BuildPagination(
				tt.page, tt.limit, tt.defaultLimit, tt.maxLimit, tt.total,
			)

			if limit != tt.expectedLimit {
				t.Errorf("expected limit %d, got %d", tt.expectedLimit, limit)
			}
			if offset != tt.expectedOffset {
				t.Errorf("expected offset %d, got %d", tt.expectedOffset, offset)
			}
			if pagination.Page != tt.expectedPage {
				t.Errorf("expected pagination.Page %d, got %d", tt.expectedPage, pagination.Page)
			}
			if pagination.Limit != tt.expectedLimit {
				t.Errorf("expected pagination.Limit %d, got %d", tt.expectedLimit, pagination.Limit)
			}
			if pagination.Total != tt.total {
				t.Errorf("expected pagination.Total %d, got %d", tt.total, pagination.Total)
			}
			if pagination.TotalPages != tt.expectedTotalPages {
				t.Errorf("expected pagination.TotalPages %d, got %d", tt.expectedTotalPages, pagination.TotalPages)
			}
		})
	}
}

func TestValidateLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		limit        int
		defaultLimit int
		maxLimit     int
		expected     int
	}{
		{
			name:         "Valid limit",
			limit:        10,
			defaultLimit: 5,
			maxLimit:     20,
			expected:     10,
		},
		{
			name:         "Limit less than or equal to 0",
			limit:        0,
			defaultLimit: 5,
			maxLimit:     20,
			expected:     5,
		},
		{
			name:         "Limit exceeds maximum",
			limit:        25,
			defaultLimit: 5,
			maxLimit:     20,
			expected:     20,
		},
		{
			name:         "Negative limit",
			limit:        -5,
			defaultLimit: 5,
			maxLimit:     20,
			expected:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := usecase.ValidateLimit(tt.limit, tt.defaultLimit, tt.maxLimit)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
