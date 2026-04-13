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

func TestGetProduct_Execute(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			GetByID(gomock.Any(), id).
			Return(testhelpers.RestoredProduct(id), nil)

		uc := v1usecase.NewGetProduct(repo)
		got, err := uc.Execute(context.Background(), v1usecase.GetProductQuery{ID: id}).ToGo()

		require.NoError(t, err)
		assert.Equal(t, id, got.Product.ID())
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			GetByID(gomock.Any(), id).
			Return(nil, apperr.NotFound("product", id))

		uc := v1usecase.NewGetProduct(repo)
		_, err := uc.Execute(context.Background(), v1usecase.GetProductQuery{ID: id}).ToGo()

		require.Error(t, err)
		assert.True(t, apperr.IsKind(err, apperr.KindNotFound))
	})

	t.Run("repository_error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			GetByID(gomock.Any(), id).
			Return(nil, errors.New("connection reset"))

		uc := v1usecase.NewGetProduct(repo)
		_, err := uc.Execute(context.Background(), v1usecase.GetProductQuery{ID: id}).ToGo()

		require.Error(t, err)
	})
}
