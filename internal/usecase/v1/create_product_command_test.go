package v1usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	apperr "github.com/selfshop-dev/lib-apperr"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/mocks"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
	"github.com/selfshop-dev/ms-catalog/testhelpers"
)

func TestCreateProduct_Execute(t *testing.T) {
	t.Parallel()

	validCmd := v1usecase.CreateProductCommand{
		Name:       "Widget",
		Slug:       "widget",
		PriceCents: 0,
		Currency:   "USD",
		Status:     domain.ProductStatusDraft,
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)

		id := uuid.New()
		repo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ *domain.Product) (*domain.Product, error) {
				// restore a product with a pre-generated ID to assert against it in the result
				return testhelpers.RestoredProduct(id), nil
			})

		uc := v1usecase.NewCreateProduct(repo)
		got, err := uc.Execute(context.Background(), validCmd).ToGo()

		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
	})

	t.Run("domain_validation_error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		// repo.Create is never called — domain rejects the command before reaching the repository
		repo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		uc := v1usecase.NewCreateProduct(repo)
		_, err := uc.Execute(context.Background(), v1usecase.CreateProductCommand{
			Name:     "", // violates domain validation
			Slug:     "widget",
			Currency: "USD",
			Status:   domain.ProductStatusDraft,
		}).ToGo()

		require.Error(t, err)
	})

	t.Run("repository_conflict", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(nil, apperr.Conflictf("product with slug widget already exists"))

		uc := v1usecase.NewCreateProduct(repo)
		_, err := uc.Execute(context.Background(), validCmd).ToGo()

		require.Error(t, err)
		assert.True(t, apperr.IsKind(err, apperr.KindConflict))
	})

	t.Run("repository_error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("connection reset"))

		uc := v1usecase.NewCreateProduct(repo)
		_, err := uc.Execute(context.Background(), validCmd).ToGo()

		require.Error(t, err)
	})
}
