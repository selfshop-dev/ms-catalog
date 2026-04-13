package v1handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
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

func TestUpdateProductHandler(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	validBody := map[string]any{
		"name":        "Updated Widget",
		"slug":        "updated-widget",
		"price_cents": 500,
		"currency":    "USD",
	}

	t.Run("200_ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateProductCommand,
			v1usecase.UpdateProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Ok[v1usecase.UpdateProductResult, error](
				v1usecase.UpdateProductResult{ID: id},
			))

		route := v1handler.NewUpdateProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String(), validBody)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data v1handler.UpdateProductResponse `json:"data"`
		}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, id, resp.Data.ID)
	})

	t.Run("400_invalid_uuid", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateProductCommand,
			v1usecase.UpdateProductResult,
		](ctrl)
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewUpdateProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/not-a-uuid", validBody)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("400_invalid_json", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateProductCommand,
			v1usecase.UpdateProductResult,
		](ctrl)
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewUpdateProduct(uc)
		r, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/products/"+id.String(), strings.NewReader("{invalid"))
		r.Header.Set("Content-Type", "application/json")
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("422_missing_required_fields", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateProductCommand,
			v1usecase.UpdateProductResult,
		](ctrl)
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewUpdateProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String(), map[string]any{
			"price_cents": 500,
			// name, slug, currency are missing
		})
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("404_not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateProductCommand,
			v1usecase.UpdateProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.UpdateProductResult{}, apperr.NotFound("product", id)))

		route := v1handler.NewUpdateProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String(), validBody)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("409_conflict", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateProductCommand,
			v1usecase.UpdateProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.UpdateProductResult{}, apperr.Conflictf("slug already exists")))

		route := v1handler.NewUpdateProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String(), validBody)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("500_internal", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateProductCommand,
			v1usecase.UpdateProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.UpdateProductResult{}, apperr.Internal("db down")))

		route := v1handler.NewUpdateProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String(), validBody)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
