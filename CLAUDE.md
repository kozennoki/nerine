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
- `BFF_API_KEY`: API key for BFF authentication
- `PORT`: Server port (default: 8080)

## API Authentication

All endpoints require `X-API-Key` header matching the `BFF_API_KEY` environment variable.

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
- category: Content reference to categories (required)
- body: Rich editor (required)
- description: Text field (required)
- image: Image field (required)

### Categories (endpoint: categories)
- name: Text field (required)

## Entity Implementation

### Article Entity
```go
type Article struct {
	ID          string
	Title       string
	Image       string
	Category    string
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