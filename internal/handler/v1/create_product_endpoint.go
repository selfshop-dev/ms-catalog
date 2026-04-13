package v1handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/handler"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
)

type CreateProductRequest struct {
	Description      *string `example:"Product Description"           json:"description,omitempty"`
	ShortDescription *string `example:"Product Short Description"     json:"short_description,omitempty" validate:"omitempty,max=256"`
	DisplayImageURL  *string `example:"https://example.com/image.jpg" json:"display_image_url,omitempty" validate:"omitempty,url"`
	Name             string  `example:"Product Name"                  json:"name"                        validate:"required,min=1,max=128"`
	Slug             string  `example:"product-name"                  json:"slug"                        validate:"required,min=1,max=128"`
	Currency         string  `example:"USD"                           json:"currency,omitempty"          validate:"omitempty,len=3"`
	Status           string  `example:"active"                        json:"status,omitempty"            validate:"omitempty,oneof=active inactive draft archived"`
	PriceCents       int64   `example:"10000"                         json:"price_cents"                 validate:"gte=0"`
}

type CreateProductResponse struct {
	ID uuid.UUID `json:"id"`
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	Creates a new product in the catalog and returns its generated UUID.
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateProductRequest	true	"Product payload"
//	@Success		201		{object}	CreateProductResponse
//	@Failure		400		{object}	response.envelope	"Invalid JSON body or failed field validation"
//	@Failure		409		{object}	response.envelope	"Slug already exists"
//	@Failure		422		{object}	response.envelope	"Semantic validation error with per-field detail"
//	@Failure		500		{object}	response.envelope	"Internal server error"
//	@Router			/products [post]
func NewCreateProduct(u usecase.Executer[
	v1usecase.CreateProductCommand,
	v1usecase.CreateProductResult,
],
) handler.Route {
	return handler.Route{
		Method:  http.MethodPost,
		Path:    "/products",
		Handler: func(w http.ResponseWriter, r *http.Request) { createProduct(u, w, r) },
	}
}

func createProduct(u usecase.Executer[
	v1usecase.CreateProductCommand,
	v1usecase.CreateProductResult,
],
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.Respond.BadRequest(w, r, "invalid json body")
		return
	}

	if err := handler.Validate(req); err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	cmd := v1usecase.CreateProductCommand{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		DisplayImageURL:  req.DisplayImageURL,
		PriceCents:       req.PriceCents,
		Currency:         req.Currency,
		Status:           domain.ProductStatus(req.Status),
	}

	res, err := u.Execute(r.Context(), cmd).ToGo()
	if err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	handler.Respond.Created(w, r, CreateProductResponse{ID: res.ID})
}
