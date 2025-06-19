package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kozennoki/nerine/internal/interfaces/middleware"
	"github.com/labstack/echo/v4"
)

func TestAPIKeyAuth_ValidAPIKey(t *testing.T) {
	t.Parallel()

	expectedAPIKey := "valid-api-key"
	middlewareFunc := middleware.APIKeyAuth(expectedAPIKey)

	// Create a mock handler
	mockHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Create Echo instance and request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", expectedAPIKey)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute middleware
	handler := middlewareFunc(mockHandler)
	err := handler(c)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got: %d", rec.Code)
	}

	if rec.Body.String() != "success" {
		t.Errorf("Expected body 'success', got: %s", rec.Body.String())
	}
}

func TestAPIKeyAuth_MissingAPIKey(t *testing.T) {
	t.Parallel()

	expectedAPIKey := "valid-api-key"
	middleware := middleware.APIKeyAuth(expectedAPIKey)

	// Create a mock handler
	mockHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Create Echo instance and request without API key
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute middleware
	handler := middleware(mockHandler)
	err := handler(c)

	if err == nil {
		t.Error("Expected error for missing API key, got nil")
	}

	// Check if it's an HTTP error with the correct status
	if httpErr, ok := err.(*echo.HTTPError); ok {
		if httpErr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got: %d", httpErr.Code)
		}
		if httpErr.Message != "missing API key" {
			t.Errorf("Expected message 'missing API key', got: %v", httpErr.Message)
		}
	} else {
		t.Errorf("Expected echo.HTTPError, got: %T", err)
	}
}

func TestAPIKeyAuth_InvalidAPIKey(t *testing.T) {
	t.Parallel()

	expectedAPIKey := "valid-api-key"
	middleware := middleware.APIKeyAuth(expectedAPIKey)

	// Create a mock handler
	mockHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Create Echo instance and request with invalid API key
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "invalid-api-key")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute middleware
	handler := middleware(mockHandler)
	err := handler(c)

	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}

	// Check if it's an HTTP error with the correct status
	if httpErr, ok := err.(*echo.HTTPError); ok {
		if httpErr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got: %d", httpErr.Code)
		}
		if httpErr.Message != "invalid API key" {
			t.Errorf("Expected message 'invalid API key', got: %v", httpErr.Message)
		}
	} else {
		t.Errorf("Expected echo.HTTPError, got: %T", err)
	}
}

func TestAPIKeyAuth_EmptyAPIKey(t *testing.T) {
	t.Parallel()

	expectedAPIKey := "valid-api-key"
	middleware := middleware.APIKeyAuth(expectedAPIKey)

	// Create a mock handler
	mockHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Create Echo instance and request with empty API key
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute middleware
	handler := middleware(mockHandler)
	err := handler(c)

	if err == nil {
		t.Error("Expected error for empty API key, got nil")
	}

	// Check if it's an HTTP error with the correct status
	if httpErr, ok := err.(*echo.HTTPError); ok {
		if httpErr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got: %d", httpErr.Code)
		}
		if httpErr.Message != "missing API key" {
			t.Errorf("Expected message 'missing API key', got: %v", httpErr.Message)
		}
	} else {
		t.Errorf("Expected echo.HTTPError, got: %T", err)
	}
}

func TestAPIKeyAuth_CaseSensitive(t *testing.T) {
	t.Parallel()

	expectedAPIKey := "Valid-API-Key"
	middleware := middleware.APIKeyAuth(expectedAPIKey)

	// Create a mock handler
	mockHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Create Echo instance and request with different case API key
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "valid-api-key") // lowercase
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute middleware
	handler := middleware(mockHandler)
	err := handler(c)

	if err == nil {
		t.Error("Expected error for case-sensitive API key mismatch, got nil")
	}

	// Check if it's an HTTP error with the correct status
	if httpErr, ok := err.(*echo.HTTPError); ok {
		if httpErr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got: %d", httpErr.Code)
		}
		if httpErr.Message != "invalid API key" {
			t.Errorf("Expected message 'invalid API key', got: %v", httpErr.Message)
		}
	} else {
		t.Errorf("Expected echo.HTTPError, got: %T", err)
	}
}
