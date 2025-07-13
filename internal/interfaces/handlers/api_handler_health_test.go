package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func TestAPIHandler_HealthCheck(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, _ := CreateTestAPIHandler(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HealthCheck(c)

	if err != nil {
		t.Errorf("HealthCheck() error = %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("HealthCheck() status = %v, want %v", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("HealthCheck() body should contain status ok, got %v", rec.Body.String())
	}
}
