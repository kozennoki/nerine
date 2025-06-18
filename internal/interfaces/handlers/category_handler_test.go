package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/usecase"
	"github.com/kozennoki/nerine/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func TestCategoryHandler_GetCategories(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetCategoriesUsecase := mocks.NewMockGetCategoriesUsecase(ctrl)
	handler := NewCategoryHandler(mockGetCategoriesUsecase)

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "正常系: カテゴリ一覧取得",
			setupMock: func() {
				categories := []*entity.Category{
					{
						Slug: "tech",
						Name: "技術",
					},
					{
						Slug: "life",
						Name: "生活",
					},
					{
						Slug: "business",
						Name: "ビジネス",
					},
				}
				output := usecase.GetCategoriesUsecaseOutput{Categories: categories}
				mockGetCategoriesUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetCategoriesUsecaseInput{}).
					Return(output, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"categories":[{"Slug":"tech","Name":"技術"},{"Slug":"life","Name":"生活"},{"Slug":"business","Name":"ビジネス"}]}`,
		},
		{
			name: "正常系: 空のカテゴリ一覧",
			setupMock: func() {
				output := usecase.GetCategoriesUsecaseOutput{Categories: []*entity.Category{}}
				mockGetCategoriesUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetCategoriesUsecaseInput{}).
					Return(output, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"categories":[]}`,
		},
		{
			name: "異常系: Usecaseエラー",
			setupMock: func() {
				e := errors.New("usecase error")
				mockGetCategoriesUsecase.EXPECT().
					Exec(gomock.Any(), usecase.GetCategoriesUsecaseInput{}).
					Return(usecase.GetCategoriesUsecaseOutput{}, e)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"detail":"usecase error","error":"Failed to get categories"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetCategories(c)

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコードが一致しません。expected: %d, got: %d", tt.expectedStatus, rec.Code)
			}

			body := strings.TrimSpace(rec.Body.String())
			if body != tt.expectedBody {
				t.Errorf("レスポンスボディが一致しません。\nexpected: %s\ngot: %s", tt.expectedBody, body)
			}
		})
	}
}
