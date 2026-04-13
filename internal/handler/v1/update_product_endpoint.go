package v1handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/selfshop-dev/ms-catalog/internal/handler"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
)

type UpdateProductRequest struct {
	Description      *string `example:"Product Description"           json:"description,omitempty"`
	ShortDescription *string `example:"Product Short Description"     json:"short_description,omitempty" validate:"omitempty,max=256"`
	DisplayImageURL  *string `example:"https://example.com/image.jpg" json:"display_image_url,omitempty" validate:"omitempty,url"`
	Name             string  `example:"Product Name"                  json:"name"                        validate:"required,min=1,max=128"`
	Slug             string  `example:"product-name"                  json:"slug"                        validate:"required,min=1,max=128"`
	Currency         string  `example:"USD"                           json:"currency"                    validate:"required,len=3"`
	PriceCents       int64   `example:"10000"                         json:"price_cents"                 validate:"gte=0"`
}

type UpdateProductResponse struct {
	ID uuid.UUID `json:"id"`
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	Fully replaces a product's fields by its UUID.
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Product UUID"	Format(uuid)
//	@Param			body	body		UpdateProductRequest	true	"Product payload"
//	@Success		200		{object}	UpdateProductResponse
//	@Failure		400		{object}	response.envelope	"Invalid UUID format or JSON body"
//	@Failure		404		{object}	response.envelope	"Product not found"
//	@Failure		409		{object}	response.envelope	"Slug already exists"
//	@Failure		422		{object}	response.envelope	"Semantic validation error with per-field detail"
//	@Failure		500		{object}	response.envelope	"Internal server error"
//	@Router			/products/{id} [put]
func NewUpdateProduct(u usecase.Executer[
	v1usecase.UpdateProductCommand,
	v1usecase.UpdateProductResult,
],
) handler.Route {
	return handler.Route{
		Method:  http.MethodPut,
		Path:    "/products/{id}",
		Handler: func(w http.ResponseWriter, r *http.Request) { updateProduct(u, w, r) },
	}
}

func updateProduct(u usecase.Executer[
	v1usecase.UpdateProductCommand,
	v1usecase.UpdateProductResult,
],
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		handler.Respond.BadRequest(w, r, "invalid id")
		return
	}

	var req UpdateProductRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.Respond.BadRequest(w, r, "invalid json body")
		return
	}

	if err = handler.Validate(req); err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	cmd := v1usecase.UpdateProductCommand{
		ID:               id,
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		DisplayImageURL:  req.DisplayImageURL,
		PriceCents:       req.PriceCents,
		Currency:         req.Currency,
	}

	res, err := u.Execute(r.Context(), cmd).ToGo()
	if err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	handler.Respond.Ok(w, r, UpdateProductResponse{ID: res.ID})
}
