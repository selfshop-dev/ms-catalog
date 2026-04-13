package v1handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/handler"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
)

type UpdateStatusProductRequest struct {
	Status string `example:"active" json:"status" validate:"required,oneof=active inactive draft archived"`
}

type UpdateStatusProductResponse struct {
	ID uuid.UUID `json:"id"`
}

// UpdateStatusProduct godoc
//
//	@Summary		Update status product
//	@Description	Changes the status of a product. Allowed values: active, inactive, draft, archived.
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Product UUID"	Format(uuid)
//	@Param			body	body		UpdateStatusProductRequest	true	"Status payload"
//	@Success		200		{object}	UpdateStatusProductResponse
//	@Failure		400		{object}	response.envelope	"Invalid UUID format or JSON body"
//	@Failure		404		{object}	response.envelope	"Product not found"
//	@Failure		422		{object}	response.envelope	"Validation error"
//	@Failure		500		{object}	response.envelope	"Internal server error"
//	@Router			/products/{id}/status [put]
func NewUpdateStatusProduct(u usecase.Executer[
	v1usecase.UpdateStatusProductCommand,
	v1usecase.UpdateStatusProductResult,
],
) handler.Route {
	return handler.Route{
		Method:  http.MethodPut,
		Path:    "/products/{id}/status",
		Handler: func(w http.ResponseWriter, r *http.Request) { updateStatusProduct(u, w, r) },
	}
}

func updateStatusProduct(u usecase.Executer[
	v1usecase.UpdateStatusProductCommand,
	v1usecase.UpdateStatusProductResult,
],
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		handler.Respond.BadRequest(w, r, "invalid id")
		return
	}

	var req UpdateStatusProductRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.Respond.BadRequest(w, r, "invalid json body")
		return
	}

	if err = handler.Validate(req); err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	cmd := v1usecase.UpdateStatusProductCommand{
		ID:     id,
		Status: domain.ProductStatus(req.Status),
	}

	res, err := u.Execute(r.Context(), cmd).ToGo()
	if err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	handler.Respond.Ok(w, r, UpdateStatusProductResponse{ID: res.ID})
}
