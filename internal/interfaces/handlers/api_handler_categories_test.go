package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func TestAPIHandler_GetCategories(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mocks := CreateTestAPIHandler(ctrl)

	tests := []struct {
		name           string
		mockOutput     usecase.GetCategoriesUsecaseOutput
		mockError      error
		expectedStatus int
	}{
		{
			name: "Success",
			mockOutput: usecase.GetCategoriesUsecaseOutput{
				Categories: []*entity.Category{
					{Slug: "tech", Name: "Technology"},
					{Slug: "business", Name: "Business"},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error from usecase",
			mockOutput:     usecase.GetCategoriesUsecaseOutput{},
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mocks.GetCategoriesUsecase.EXPECT().
				Exec(gomock.Any(), usecase.GetCategoriesUsecaseInput{}).
				Return(tt.mockOutput, tt.mockError)

			err := handler.GetCategories(c)

			if err != nil {
				t.Errorf("GetCategories() error = %v", err)
			}
			if rec.Code != tt.expectedStatus {
				t.Errorf("GetCategories() status = %v, want %v", rec.Code, tt.expectedStatus)
			}
		})
	}
}
