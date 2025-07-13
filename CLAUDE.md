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
task test              # Run all tests
go test ./...          # Alternative command
go test -v ./...       # Verbose output

# Generate mocks
task generate-mocks    # Generate all mocks
task gen-mocks         # Alias for generate-mocks

# Generate OpenAPI code
task generate-openapi  # Generate Go code from OpenAPI schema
task gen-api          # Alias for generate-openapi

# Test coverage
task coverage          # Generate coverage report
task cov               # Alias for coverage
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

## Testing Guidelines

### Test Structure
- Use `package xxx_test` pattern for external testing
- Create `export_test.go` files to expose private functions/methods for testing
- Organize tests with clear naming: `TestFunctionName_Scenario`

### Test Execution
- Enable parallel execution with `t.Parallel()` for tests that don't share global state
- Avoid parallel execution for tests that modify environment variables or global state
- Use table-driven tests for multiple scenarios

### Mock Generation and Usage
- Generate mocks using `task generate-mocks` or `go.uber.org/mock/mockgen`
- Store mocks in `mocks/` directories within each layer
- Use dependency injection to make code testable with mocks

### Test Coverage
- Target 80%+ test coverage for production code
- Focus on business logic in UseCase and Handler layers
- Test both success and error scenarios
- Use `task coverage` to generate coverage reports

### Test Organization by Layer

#### UseCase Layer Testing
```go
// Test pattern with repository mocks
func TestUseCaseName_Exec(t *testing.T) {
    t.Parallel()

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockRepository(ctrl)
    usecase := NewUseCase(mockRepo)

    // Setup expectations and test
}
```

#### Handler Layer Testing
```go
// Test pattern with usecase mocks
func TestHandlerName_Method(t *testing.T) {
    t.Parallel()

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUseCase := mocks.NewMockUseCase(ctrl)
    handler := NewHandler(mockUseCase)

    // Setup Echo context and test HTTP handling
}
```

#### Infrastructure Layer Testing
```go
// Test pattern for configuration and utilities
func TestConfigLoad_Success(t *testing.T) {
    // No t.Parallel() for environment variable tests

    os.Setenv("KEY", "value")
    defer os.Unsetenv("KEY")

    // Test configuration loading
}

func TestUtilityFunction(t *testing.T) {
    t.Parallel() // Safe for pure functions

    // Test utility functions
}
```

### Environment Variable Testing
- Tests modifying environment variables should NOT use `t.Parallel()`
- Use `defer os.Unsetenv()` to clean up after tests
- Group environment variable tests to run sequentially

### Best Practices
- Write tests before or alongside implementation (TDD/BDD approach)
- Test error conditions and edge cases thoroughly
- Use meaningful test names that describe the scenario
- Keep tests focused and independent
- Use `export_test.go` for testing private functions when necessary

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

## OpenAPI Schema-Driven Development

This project uses OpenAPI schema-driven development with automatic code generation.

### Schema Repository
- **Repository**: `hibiscus` (https://github.com/kozennoki/hibiscus.git)
- **Local Path**: `schema/`
- **Main File**: `schema/openapi.yaml`

### Schema Management
```bash
# Update schema to latest version
git submodule update --remote schema
git add schema
git commit -m "📚 API スキーマを更新"

# Initialize submodules (for new clones)
git submodule update --init --recursive
```

### Code Generation Setup

The project uses `oapi-codegen v2` for automatic Go code generation from OpenAPI specifications.

**Generated Code Location**: `internal/openapi/openapi.go`

**Generated Contents**:
- Type definitions (Article, Category, Pagination, etc.)
- ServerInterface for all API endpoints
- Echo handlers with parameter parsing and validation
- Response models matching OpenAPI specification

### Integration Pattern

1. **Handlers implement ServerInterface**: All handlers implement `openapi.ServerInterface`
2. **Type conversion via Presenter**: `internal/interfaces/presenter/converter.go` converts between domain entities and generated types
3. **Automatic routing**: Generated routes are registered via `openapi.RegisterHandlers()`

### Development Workflow for API Changes

When making API changes:
1. **Update OpenAPI schema** in `schema/openapi.yaml`
2. **Regenerate code** with `task gen-api`
3. **Update converters** in `presenter/converter.go` if new fields are added
4. **Test integration** to ensure type compatibility

### Benefits

- **Schema-first design**: API specification drives implementation
- **Type safety**: Generated types ensure API compliance
- **Automatic validation**: Parameter parsing and validation generated from schema
- **Documentation sync**: Code always matches specification

The OpenAPI schema serves as the single source of truth for API specifications shared between:
- **Nerine** (BFF API - this repository)
- **Abelia** (Frontend Next.js application)

## Development Workflow

When starting work on this project, follow these steps:
1. **Read README.md** to understand the project structure and design
2. **Check docs/tasks.md** to see current tasks and priorities
3. **Update docs/tasks.md** after completing each task to maintain progress tracking
4. **Use schema-driven development**: When making API changes, update OpenAPI schema first, then regenerate code with `task gen-api`
5. Follow Clean Architecture principles and maintain layer separation throughout implementation
