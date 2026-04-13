package v1handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestCreateProductHandler(t *testing.T) {
	t.Parallel()

	validBody := map[string]any{
		"name":        "Widget",
		"slug":        "widget",
		"price_cents": 1000,
		"currency":    "USD",
		"status":      "draft",
	}

	t.Run("201_created", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.CreateProductCommand,
			v1usecase.CreateProductResult,
		](ctrl)

		id := uuid.New()
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Ok[v1usecase.CreateProductResult, error](
				v1usecase.CreateProductResult{ID: id},
			))

		route := v1handler.NewCreateProduct(uc)
		w := httptest.NewRecorder()
		route.Handler(w, testhelpers.NewJSONRequest(t, http.MethodPost, "/products", validBody))

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp struct {
			Data v1handler.CreateProductResponse `json:"data"`
		}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, id, resp.Data.ID)
	})

	t.Run("400_invalid_json", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.CreateProductCommand,
			v1usecase.CreateProductResult,
		](ctrl)
		// Execute is never called — handler rejects malformed JSON before reaching usecase
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewCreateProduct(uc)
		w := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/products", bytes.NewBufferString("{invalid"))
		r.Header.Set("Content-Type", "application/json")
		route.Handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	})

	t.Run("422_missing_required_fields", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.CreateProductCommand,
			v1usecase.CreateProductResult,
		](ctrl)
		// Execute is never called — handler rejects invalid request before reaching usecase
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewCreateProduct(uc)
		w := httptest.NewRecorder()
		route.Handler(w, testhelpers.NewJSONRequest(t, http.MethodPost, "/products", map[string]any{
			"price_cents": 1000,
			// name and slug are missing
		}))

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	})

	t.Run("409_conflict", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.CreateProductCommand,
			v1usecase.CreateProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(
				v1usecase.CreateProductResult{},
				apperr.Conflictf("product with slug widget already exists"),
			))

		route := v1handler.NewCreateProduct(uc)
		w := httptest.NewRecorder()
		route.Handler(w, testhelpers.NewJSONRequest(t, http.MethodPost, "/products", validBody))

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("422_domain_validation", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.CreateProductCommand,
			v1usecase.CreateProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(
				v1usecase.CreateProductResult{},
				apperr.Unprocessable("invalid product"),
			))

		route := v1handler.NewCreateProduct(uc)
		w := httptest.NewRecorder()
		route.Handler(w, testhelpers.NewJSONRequest(t, http.MethodPost, "/products", validBody))

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("500_internal", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.CreateProductCommand,
			v1usecase.CreateProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(
				v1usecase.CreateProductResult{},
				apperr.Internal("db down"),
			))

		route := v1handler.NewCreateProduct(uc)
		w := httptest.NewRecorder()
		route.Handler(w, testhelpers.NewJSONRequest(t, http.MethodPost, "/products", validBody))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
