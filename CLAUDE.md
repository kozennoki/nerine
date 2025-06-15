# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Nerine is a BFF (Backend for Frontend) API built with Go for a Next.js + microCMS blog system. The project follows DDD (Domain-Driven Design) principles with Clean Architecture patterns.

## Architecture

- **Domain Layer**: `internal/domain/` - Contains business rules and core logic
  - `entity/`: Domain entities (Article, Category)
  - `service/`: Domain services for complex business logic
  - `repository/`: Repository interfaces for data access
- **Use Case Layer**: `internal/usecase/` - Application services implementing business logic with DTOs
- **Infrastructure Layer**: `internal/infrastructure/` - External dependencies (microCMS SDK, logger)
- **Interface Layer**: `internal/interfaces/` - HTTP handlers, middleware, and presenters using Echo framework

## Development Commands

```bash
# Initialize Go module (if not done)
go mod init github.com/kozennoki/nerine

# Install dependencies
go mod tidy

# Run development server
go run cmd/server/main.go

# Run tests
go test ./...
```

## Environment Variables

Required environment variables:
- `MICROCMS_API_KEY`: microCMS API key
- `MICROCMS_SERVICE_ID`: microCMS service ID
- `NERINE_API_KEY`: API key for nerine authentication
- `PORT`: Server port (default: 8080)

## API Authentication

All endpoints require `X-API-Key` header matching the `NERINE_API_KEY` environment variable.

## Key Dependencies

- Echo v4: HTTP framework
- zap: Structured logging
- microCMS Go SDK: CMS integration
- gomock: Mock generation for unit testing
- Standard library for configuration (`os.Getenv`)

## Development Principles

- Use simple, standard library solutions first before adding external dependencies
- Follow Clean Architecture with clear separation between layers
- Implement DDD patterns with domain-specific packages
- Use DTOs (Data Transfer Objects) for layer communication to maintain decoupling
- Implement dependency injection for testability
- Generate mocks with gomock for unit testing
- Structure logging with zap for production readiness

## microCMS API Schema

### Blog (endpoint: blog)
- title: Text field (required)
- category: Content reference to categories (required, returns object with id and name)
- body: Rich editor (required)
- description: Text field (required)
- image: Image field (required, returns object with url, height, width)

### Categories (endpoint: categories)
- name: Text field (required)

## microCMS Response Structure

### Article Field
```go
type article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Image       image     `json:"image"`
	Category    category  `json:"category"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
```

### Image Field
```go
type image struct {
    URL    string `json:"url"`
    Height int    `json:"height"`
    Width  int    `json:"width"`
}
```

### Category Reference
```go
type category struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

## Entity Implementation

### Article Entity
```go
type Article struct {
	ID          string
	Title       string
	Image       string
	Category    Category
	Description string
	Body        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
```

### Category Entity
```go
type Category struct {
	Slug string
	Name string
}
```

## Data Mapping

### microCMS to Entity Conversion

**Article Mapping:**
- `item.Image.URL` → `entity.Article.Image` (string)
- `item.Category.ID` → `entity.Article.Category.Slug` (string)
- `item.Category.Name` → `entity.Article.Category.Name` (string)

**Key Points:**
- microCMS image field is an object containing URL, height, width - we extract only the URL
- microCMS category field is a reference object containing ID and name - we map ID to Slug and name to Name
- All other fields are mapped directly without transformation

## UseCase Pattern

Each usecase follows this structure:
- Interface definition with `Exec` method
- Input/Output DTOs embedded in the same file
- Constructor for dependency injection
- Implementation struct with repository dependencies

Example:
```go
type GetArticlesUsecase interface {
    Exec(context.Context, GetArticlesUsecaseInput) (GetArticlesUsecaseOutput, error)
}

func NewGetArticles(
    repo repository.ArticleRepository,
) GetArticlesUsecase {
    return &getArticles{
        repo: repo,
    }
}
```

## Development Workflow

When starting work on this project, follow these steps:
1. **Read README.md** to understand the project structure and design
2. **Check docs/tasks.md** to see current tasks and priorities
3. **Update docs/tasks.md** after completing each task to maintain progress tracking
4. Follow Clean Architecture principles and maintain layer separation throughout implementation
