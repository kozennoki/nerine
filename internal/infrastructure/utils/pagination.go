package utils

// Pagination represents pagination metadata for API responses
type Pagination struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"totalPages"`
}

// ConvertPageToOffset converts page-based pagination to offset-based pagination
// page: 1-based page number
// limit: number of items per page
// Returns: offset for database queries
func ConvertPageToOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}

// CalculateTotalPages calculates total pages from total count and limit
func CalculateTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	if total <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}

// NewPagination creates a new Pagination struct
func NewPagination(total, page, limit int) Pagination {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalPages := CalculateTotalPages(total, limit)

	return Pagination{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
