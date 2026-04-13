package v1handler_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	apperr "github.com/selfshop-dev/lib-apperr"
	result "github.com/selfshop-dev/lib-result"

	v1handler "github.com/selfshop-dev/ms-catalog/internal/handler/v1"
	"github.com/selfshop-dev/ms-catalog/internal/mocks"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
	"github.com/selfshop-dev/ms-catalog/testhelpers"
)

func TestDeleteProductHandler(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	t.Run("204_no_content", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.DeleteProductCommand,
			v1usecase.DeleteProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), v1usecase.DeleteProductCommand{ID: id}).
			Return(result.Ok[v1usecase.DeleteProductResult, error](v1usecase.DeleteProductResult{}))

		route := v1handler.NewDeleteProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodDelete, "/products/"+id.String(), nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("400_invalid_uuid", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.DeleteProductCommand,
			v1usecase.DeleteProductResult,
		](ctrl)
		// Execute is never called — handler rejects malformed UUID before reaching usecase
		uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(0)

		route := v1handler.NewDeleteProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodDelete, "/products/not-a-uuid", nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("404_not_found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.DeleteProductCommand,
			v1usecase.DeleteProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.DeleteProductResult{}, apperr.NotFound("product", id)))

		route := v1handler.NewDeleteProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodDelete, "/products/"+id.String(), nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("500_internal", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		uc := mocks.NewMockExecuter[
			v1usecase.DeleteProductCommand,
			v1usecase.DeleteProductResult,
		](ctrl)
		uc.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(result.Of(v1usecase.DeleteProductResult{}, apperr.Internal("db down")))

		route := v1handler.NewDeleteProduct(uc)
		r := testhelpers.NewJSONRequest(t, http.MethodDelete, "/products/"+id.String(), nil)
		w := testhelpers.ServeWithChi(t, route, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
