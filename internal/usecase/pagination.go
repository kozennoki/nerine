package usecase

import "github.com/kozennoki/nerine/internal/infrastructure/utils"

// BuildPagination validates page and limit parameters and builds pagination info
// Returns: validatedLimit, offset, pagination
func BuildPagination(page, limit, defaultLimit, maxLimit int, total int) (int, int, utils.Pagination) {
	validatedPage := page
	if validatedPage < 1 {
		validatedPage = 1
	}

	validatedLimit := limit
	if validatedLimit <= 0 {
		validatedLimit = defaultLimit
	}
	if validatedLimit > maxLimit {
		validatedLimit = maxLimit
	}

	offset := utils.ConvertPageToOffset(validatedPage, validatedLimit)
	pagination := utils.NewPagination(total, validatedPage, validatedLimit)

	return validatedLimit, offset, pagination
}

// ValidateLimit validates limit parameter with default and maximum values
func ValidateLimit(limit, defaultLimit, maxLimit int) int {
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}