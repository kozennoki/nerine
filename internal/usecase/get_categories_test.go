package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/kozennoki/nerine/internal/domain/entity"
	"github.com/kozennoki/nerine/internal/domain/repository/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetCategories_Exec(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCategoryRepository(ctrl)

	tests := []struct {
		name          string
		input         GetCategoriesUsecaseInput
		setupMock     func()
		expected      GetCategoriesUsecaseOutput
		expectedError error
	}{
		{
			name:  "正常系: カテゴリ一覧が正常に取得できる",
			input: GetCategoriesUsecaseInput{},
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
				mockRepo.EXPECT().GetCategories(gomock.Any()).Return(categories, nil)
			},
			expected: GetCategoriesUsecaseOutput{
				Categories: []*entity.Category{
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
				},
			},
			expectedError: nil,
		},
		{
			name:  "正常系: 空のカテゴリ一覧",
			input: GetCategoriesUsecaseInput{},
			setupMock: func() {
				mockRepo.EXPECT().GetCategories(gomock.Any()).Return([]*entity.Category{}, nil)
			},
			expected: GetCategoriesUsecaseOutput{
				Categories: []*entity.Category{},
			},
			expectedError: nil,
		},
		{
			name:  "異常系: リポジトリエラー",
			input: GetCategoriesUsecaseInput{},
			setupMock: func() {
				e := errors.New("repository error")
				mockRepo.EXPECT().GetCategories(gomock.Any()).Return(nil, e)
			},
			expected:      GetCategoriesUsecaseOutput{},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()

			usecase := NewGetCategories(mockRepo)
			result, err := usecase.Exec(context.Background(), tt.input)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("期待されるエラーが発生しませんでした。expected: %v", tt.expectedError)
					return
				}
				if errors.Is(err, tt.expectedError) {
					t.Errorf("エラーメッセージが一致しません。expected: %v, got: %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if len(result.Categories) != len(tt.expected.Categories) {
				t.Errorf("カテゴリ数が一致しません。expected: %d, got: %d", len(tt.expected.Categories), len(result.Categories))
				return
			}

			for i, category := range result.Categories {
				expected := tt.expected.Categories[i]
				if category.Slug != expected.Slug {
					t.Errorf("カテゴリSlug[%d]が一致しません。expected: %s, got: %s", i, expected.Slug, category.Slug)
				}
				if category.Name != expected.Name {
					t.Errorf("カテゴリ名[%d]が一致しません。expected: %s, got: %s", i, expected.Name, category.Name)
				}
			}
		})
	}
}
