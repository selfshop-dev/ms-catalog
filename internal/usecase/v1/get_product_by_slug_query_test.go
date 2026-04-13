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

func TestGetProductBySlug_Execute(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)

		id := uuid.New()
		repo.EXPECT().
			GetBySlug(gomock.Any(), "widget").
			Return(testhelpers.RestoredProduct(id), nil)

		uc := v1usecase.NewGetProductBySlug(repo)
		got, err := uc.Execute(context.Background(), v1usecase.GetProductBySlugQuery{Slug: "widget"}).ToGo()

		require.NoError(t, err)
		assert.Equal(t, id, got.Product.ID())
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			GetBySlug(gomock.Any(), "no-such-slug").
			Return(nil, apperr.NotFound("product", "no-such-slug"))

		uc := v1usecase.NewGetProductBySlug(repo)
		_, err := uc.Execute(context.Background(), v1usecase.GetProductBySlugQuery{Slug: "no-such-slug"}).ToGo()

		require.Error(t, err)
		assert.True(t, apperr.IsKind(err, apperr.KindNotFound))
	})

	t.Run("repository_error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockProductRepository(ctrl)
		repo.EXPECT().
			GetBySlug(gomock.Any(), "widget").
			Return(nil, errors.New("connection reset"))

		uc := v1usecase.NewGetProductBySlug(repo)
		_, err := uc.Execute(context.Background(), v1usecase.GetProductBySlugQuery{Slug: "widget"}).ToGo()

		require.Error(t, err)
	})
}
