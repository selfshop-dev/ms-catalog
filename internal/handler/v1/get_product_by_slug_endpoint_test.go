package v1handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	apperr "github.com/selfshop-dev/lib-apperr"
	result "github.com/selfshop-dev/lib-result"

	v1handler "github.com/selfshop-dev/ms-catalog/internal/handler/v1"
	"github.com/selfshop-dev/ms-catalog/internal/mocks"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
	"github.com/selfshop-dev/ms-catalog/testhelpers"
)

func TestGetProductBySlugHandler(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	t.Run("200_ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetProductBySlugQuery,
			v1usecase.GetProductBySlugResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), v1usecase.GetProductBySlugQuery{Slug: "widget"}).
			Return(result.Ok[v1usecase.GetProductBySlugResult, error](
				v1usecase.GetProductBySlugResult{Product: testhelpers.RestoredProduct(id)},
			))

		route := v1handler.NewGetProductBySlug(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/widget/by-slug", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data v1handler.GetProductResponse `json:"data"`
		}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, id, resp.Data.ID)
	})

	t.Run("404_not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetProductBySlugQuery,
			v1usecase.GetProductBySlugResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.GetProductBySlugResult{}, apperr.NotFound("product", "no-such-slug")))

		route := v1handler.NewGetProductBySlug(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/no-such-slug/by-slug", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("500_internal", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetProductBySlugQuery,
			v1usecase.GetProductBySlugResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.GetProductBySlugResult{}, apperr.Internal("db down")))

		route := v1handler.NewGetProductBySlug(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/widget/by-slug", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
