package v1usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/mocks"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
	"github.com/selfshop-dev/ms-catalog/testhelpers"
)

func TestGetActiveProducts_Execute(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)

		products := []*domain.Product{
			testhelpers.RestoredProduct(uuid.New()),
			testhelpers.RestoredProduct(uuid.New()),
		}
		repo.EXPECT().
			GetListActiveProducts(gomock.Any(), int32(20), int32(0)).
			Return(products, nil)

		uc := v1usecase.NewGetActiveProducts(repo)
		got, err := uc.Execute(context.Background(), v1usecase.GetActiveProductsQuery{Limit: 20, Offset: 0}).ToGo()

		require.NoError(t, err)
		assert.Len(t, got.Products, 2)
	})

	t.Run("empty_result", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			GetListActiveProducts(gomock.Any(), int32(20), int32(0)).
			Return([]*domain.Product{}, nil)

		uc := v1usecase.NewGetActiveProducts(repo)
		got, err := uc.Execute(context.Background(), v1usecase.GetActiveProductsQuery{}).ToGo()

		require.NoError(t, err)
		assert.Empty(t, got.Products)
	})

	t.Run("limit_zero_defaults_to_20", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		// limit=0 must be clamped to 20 by the usecase
		repo.EXPECT().
			GetListActiveProducts(gomock.Any(), int32(20), int32(0)).
			Return([]*domain.Product{}, nil)

		uc := v1usecase.NewGetActiveProducts(repo)
		_, err := uc.Execute(context.Background(), v1usecase.GetActiveProductsQuery{Limit: 0}).ToGo()

		require.NoError(t, err)
	})

	t.Run("limit_exceeds_100_defaults_to_20", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		// limit=200 must be clamped to 20 by the usecase
		repo.EXPECT().
			GetListActiveProducts(gomock.Any(), int32(20), int32(0)).
			Return([]*domain.Product{}, nil)

		uc := v1usecase.NewGetActiveProducts(repo)
		_, err := uc.Execute(context.Background(), v1usecase.GetActiveProductsQuery{Limit: 200}).ToGo()

		require.NoError(t, err)
	})

	t.Run("negative_offset_clamped_to_zero", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		// offset=-5 must be clamped to 0 by the usecase
		repo.EXPECT().
			GetListActiveProducts(gomock.Any(), int32(20), int32(0)).
			Return([]*domain.Product{}, nil)

		uc := v1usecase.NewGetActiveProducts(repo)
		_, err := uc.Execute(context.Background(), v1usecase.GetActiveProductsQuery{Offset: -5}).ToGo()

		require.NoError(t, err)
	})

	t.Run("repository_error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			GetListActiveProducts(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("connection reset"))

		uc := v1usecase.NewGetActiveProducts(repo)
		_, err := uc.Execute(context.Background(), v1usecase.GetActiveProductsQuery{Limit: 20}).ToGo()

		require.Error(t, err)
	})
}
