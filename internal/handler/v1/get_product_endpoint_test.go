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

func TestGetProductHandler(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	t.Run("200_ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetProductQuery,
			v1usecase.GetProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), v1usecase.GetProductQuery{ID: id}).
			Return(result.Ok[v1usecase.GetProductResult, error](
				v1usecase.GetProductResult{Product: testhelpers.RestoredProduct(id)},
			))

		route := v1handler.NewGetProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/"+id.String(), nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp struct {
			Data v1handler.GetProductResponse `json:"data"`
		}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, id, resp.Data.ID)
	})

	t.Run("400_invalid_uuid", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetProductQuery,
			v1usecase.GetProductResult,
		](ctrl)
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewGetProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/not-a-uuid", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("404_not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetProductQuery,
			v1usecase.GetProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.GetProductResult{}, apperr.NotFound("product", id)))

		route := v1handler.NewGetProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/"+id.String(), nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("500_internal", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.GetProductQuery,
			v1usecase.GetProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.GetProductResult{}, apperr.Internal("db down")))

		route := v1handler.NewGetProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodGet, "/products/"+id.String(), nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
