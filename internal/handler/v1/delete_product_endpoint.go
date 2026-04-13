package v1handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/selfshop-dev/ms-catalog/internal/handler"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
)

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Permanently deletes a product by its UUID.
//	@Tags			products
//	@Param			id	path	string	true	"Product UUID"	Format(uuid)
//	@Success		204	"No Content"
//	@Failure		400	{object}	response.envelope	"Invalid UUID format"
//	@Failure		404	{object}	response.envelope	"Product not found"
//	@Failure		500	{object}	response.envelope	"Internal server error"
//	@Router			/products/{id} [delete]
func NewDeleteProduct(u usecase.Executer[
	v1usecase.DeleteProductCommand,
	v1usecase.DeleteProductResult,
],
) handler.Route {
	return handler.Route{
		Method:  http.MethodDelete,
		Path:    "/products/{id}",
		Handler: func(w http.ResponseWriter, r *http.Request) { deleteProduct(u, w, r) },
	}
}

func deleteProduct(u usecase.Executer[
	v1usecase.DeleteProductCommand,
	v1usecase.DeleteProductResult,
],
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		handler.Respond.BadRequest(w, r, "invalid id")
		return
	}

	cmd := v1usecase.DeleteProductCommand{ID: id}

	_, err = u.Execute(r.Context(), cmd).ToGo()
	if err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	handler.Respond.NoContent(w, r)
}
