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

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	v1handler "github.com/selfshop-dev/ms-catalog/internal/handler/v1"
	"github.com/selfshop-dev/ms-catalog/internal/mocks"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
	"github.com/selfshop-dev/ms-catalog/testhelpers"
)

func TestGetActiveProductsHandler(t *testing.T) {
	t.Parallel()

	t.Run("200_ok_default_pagination", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetActiveProductsQuery,
			v1usecase.GetActiveProductsResult,
		](ctrl)

		products := []*domain.Product{
			testhelpers.RestoredProduct(uuid.New()),
			testhelpers.RestoredProduct(uuid.New()),
		}
		uc.EXPECT().
			Execute(gomock.Any(), v1usecase.GetActiveProductsQuery{Limit: 20, Offset: 0}).
			Return(result.Ok[v1usecase.GetActiveProductsResult, error](
				v1usecase.GetActiveProductsResult{Products: products},
			))

		route := v1handler.NewGetActiveProducts(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/active", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data v1handler.GetActiveProductsResponse `json:"data"`
		}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Len(t, resp.Data.Items, 2)
		assert.Equal(t, int32(20), resp.Data.Limit)
		assert.Equal(t, int32(0), resp.Data.Offset)
	})

	t.Run("200_ok_custom_pagination", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetActiveProductsQuery,
			v1usecase.GetActiveProductsResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), v1usecase.GetActiveProductsQuery{Limit: 10, Offset: 5}).
			Return(result.Ok[v1usecase.GetActiveProductsResult, error](
				v1usecase.GetActiveProductsResult{Products: []*domain.Product{}},
			))

		route := v1handler.NewGetActiveProducts(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/active?limit=10&offset=5", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("200_ok_invalid_limit_falls_back_to_default", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetActiveProductsQuery,
			v1usecase.GetActiveProductsResult,
		](ctrl)
		// invalid limit is silently ignored — handler falls back to 20
		uc.EXPECT().
			Execute(gomock.Any(), v1usecase.GetActiveProductsQuery{Limit: 20, Offset: 0}).
			Return(result.Ok[v1usecase.GetActiveProductsResult, error](
				v1usecase.GetActiveProductsResult{Products: []*domain.Product{}},
			))

		route := v1handler.NewGetActiveProducts(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/active?limit=999", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("500_internal", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetActiveProductsQuery,
			v1usecase.GetActiveProductsResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.GetActiveProductsResult{}, apperr.Internal("db down")))

		route := v1handler.NewGetActiveProducts(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/active", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
