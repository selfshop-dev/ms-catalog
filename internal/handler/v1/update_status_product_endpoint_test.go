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

func TestUpdateStatusProductHandler(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	t.Run("200_ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateStatusProductCommand,
			v1usecase.UpdateStatusProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Ok[v1usecase.UpdateStatusProductResult, error](
				v1usecase.UpdateStatusProductResult{ID: id},
			))

		route := v1handler.NewUpdateStatusProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String()+"/status",
			map[string]any{"status": "active"},
		)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data v1handler.UpdateStatusProductResponse `json:"data"`
		}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, id, resp.Data.ID)
	})

	t.Run("400_invalid_uuid", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateStatusProductCommand,
			v1usecase.UpdateStatusProductResult,
		](ctrl)
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewUpdateStatusProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/not-a-uuid/status",
			map[string]any{"status": "active"},
		)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("400_invalid_json", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateStatusProductCommand,
			v1usecase.UpdateStatusProductResult,
		](ctrl)
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewUpdateStatusProduct(uc)
		r, _ := http.NewRequestWithContext(context.Background(), http.MethodPut,
			"/products/"+id.String()+"/status",
			strings.NewReader("{invalid"),
		)
		r.Header.Set("Content-Type", "application/json")
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("422_invalid_status", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateStatusProductCommand,
			v1usecase.UpdateStatusProductResult,
		](ctrl)
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewUpdateStatusProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String()+"/status",
			map[string]any{"status": "unknown-status"},
		)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("404_not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateStatusProductCommand,
			v1usecase.UpdateStatusProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.UpdateStatusProductResult{}, apperr.NotFound("product", id)))

		route := v1handler.NewUpdateStatusProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String()+"/status",
			map[string]any{"status": "active"},
		)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("500_internal", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.UpdateStatusProductCommand,
			v1usecase.UpdateStatusProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.UpdateStatusProductResult{}, apperr.Internal("db down")))

		route := v1handler.NewUpdateStatusProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodPut, "/products/"+id.String()+"/status",
			map[string]any{"status": "active"},
		)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
