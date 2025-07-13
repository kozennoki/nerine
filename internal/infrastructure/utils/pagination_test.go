package utils_test

import (
	"testing"

	"github.com/kozennoki/nerine/internal/infrastructure/utils"
)

func TestConvertPageToOffset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		page  int
		limit int
		want  int
	}{
		{
			name:  "first page",
			page:  1,
			limit: 10,
			want:  0,
		},
		{
			name:  "second page",
			page:  2,
			limit: 10,
			want:  10,
		},
		{
			name:  "third page",
			page:  3,
			limit: 5,
			want:  10,
		},
		{
			name:  "invalid page (less than 1)",
			page:  0,
			limit: 10,
			want:  0,
		},
		{
			name:  "negative page",
			page:  -1,
			limit: 10,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := utils.ConvertPageToOffset(tt.page, tt.limit)
			if got != tt.want {
				t.Errorf("ConvertPageToOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateTotalPages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		total int
		limit int
		want  int
	}{
		{
			name:  "exact division",
			total: 20,
			limit: 10,
			want:  2,
		},
		{
			name:  "with remainder",
			total: 25,
			limit: 10,
			want:  3,
		},
		{
			name:  "single page",
			total: 5,
			limit: 10,
			want:  1,
		},
		{
			name:  "zero total",
			total: 0,
			limit: 10,
			want:  0,
		},
		{
			name:  "zero limit",
			total: 10,
			limit: 0,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := utils.CalculateTotalPages(tt.total, tt.limit)
			if got != tt.want {
				t.Errorf("CalculateTotalPages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPagination(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		total     int
		page      int
		limit     int
		wantTotal int
		wantPage  int
		wantLimit int
		wantPages int
	}{
		{
			name:      "normal case",
			total:     25,
			page:      2,
			limit:     10,
			wantTotal: 25,
			wantPage:  2,
			wantLimit: 10,
			wantPages: 3,
		},
		{
			name:      "invalid page",
			total:     25,
			page:      0,
			limit:     10,
			wantTotal: 25,
			wantPage:  1,
			wantLimit: 10,
			wantPages: 3,
		},
		{
			name:      "invalid limit",
			total:     25,
			page:      2,
			limit:     0,
			wantTotal: 25,
			wantPage:  2,
			wantLimit: 10,
			wantPages: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := utils.NewPagination(tt.total, tt.page, tt.limit)

			if got.Total != tt.wantTotal {
				t.Errorf("NewPagination().Total = %v, want %v", got.Total, tt.wantTotal)
			}
			if got.Page != tt.wantPage {
				t.Errorf("NewPagination().Page = %v, want %v", got.Page, tt.wantPage)
			}
			if got.Limit != tt.wantLimit {
				t.Errorf("NewPagination().Limit = %v, want %v", got.Limit, tt.wantLimit)
			}
			if got.TotalPages != tt.wantPages {
				t.Errorf("NewPagination().TotalPages = %v, want %v", got.TotalPages, tt.wantPages)
			}
		})
	}
}
