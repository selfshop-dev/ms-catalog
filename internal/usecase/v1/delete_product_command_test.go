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
)

func TestDeleteProduct_Execute(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		uc := v1usecase.NewDeleteProduct(repo)
		_, err := uc.Execute(context.Background(), v1usecase.DeleteProductCommand{ID: id}).ToGo()

		require.NoError(t, err)
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().Delete(gomock.Any(), id).Return(apperr.NotFound("product", id))

		uc := v1usecase.NewDeleteProduct(repo)
		_, err := uc.Execute(context.Background(), v1usecase.DeleteProductCommand{ID: id}).ToGo()

		require.Error(t, err)
		assert.True(t, apperr.IsKind(err, apperr.KindNotFound))
	})

	t.Run("repository_error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("connection reset"))

		uc := v1usecase.NewDeleteProduct(repo)
		_, err := uc.Execute(context.Background(), v1usecase.DeleteProductCommand{ID: id}).ToGo()

		require.Error(t, err)
	})
}
