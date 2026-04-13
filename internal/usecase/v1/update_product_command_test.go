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

	"github.com/selfshop-dev/ms-catalog/internal/mocks"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
	"github.com/selfshop-dev/ms-catalog/testhelpers"
)

func TestUpdateProduct_Execute(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	validCmd := v1usecase.UpdateProductCommand{
		ID:         id,
		Name:       "Updated Widget",
		Slug:       "updated-widget",
		PriceCents: 500,
		Currency:   "USD",
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			Update(gomock.Any(), id, gomock.Any()).
			Return(testhelpers.RestoredProduct(id), nil)

		uc := v1usecase.NewUpdateProduct(repo)
		got, err := uc.Execute(context.Background(), validCmd).ToGo()

		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
	})

	t.Run("domain_validation_error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		// Update is never called — domain rejects the command before reaching the repository
		repo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		uc := v1usecase.NewUpdateProduct(repo)
		_, err := uc.Execute(context.Background(), v1usecase.UpdateProductCommand{
			ID:       id,
			Name:     "", // violates domain validation
			Slug:     "updated-widget",
			Currency: "USD",
		}).ToGo()

		require.Error(t, err)
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			Update(gomock.Any(), id, gomock.Any()).
			Return(nil, apperr.NotFound("product", id))

		uc := v1usecase.NewUpdateProduct(repo)
		_, err := uc.Execute(context.Background(), validCmd).ToGo()

		require.Error(t, err)
		assert.True(t, apperr.IsKind(err, apperr.KindNotFound))
	})

	t.Run("repository_conflict", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			Update(gomock.Any(), id, gomock.Any()).
			Return(nil, apperr.Conflictf("product with slug updated-widget already exists"))

		uc := v1usecase.NewUpdateProduct(repo)
		_, err := uc.Execute(context.Background(), validCmd).ToGo()

		require.Error(t, err)
		assert.True(t, apperr.IsKind(err, apperr.KindConflict))
	})

	t.Run("repository_error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			Update(gomock.Any(), id, gomock.Any()).
			Return(nil, errors.New("connection reset"))

		uc := v1usecase.NewUpdateProduct(repo)
		_, err := uc.Execute(context.Background(), validCmd).ToGo()

		require.Error(t, err)
	})
}
