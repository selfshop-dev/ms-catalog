package v1handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/selfshop-dev/ms-catalog/internal/handler"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
)

type GetProductResponse struct {
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Description      *string   `json:"description,omitempty"`
	ShortDescription *string   `json:"short_description,omitempty"`
	DisplayImageURL  *string   `json:"display_image_url,omitempty"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Currency         string    `json:"currency"`
	Status           string    `json:"status"`
	PriceCents       int64     `json:"price_cents"`
	ID               uuid.UUID `json:"id"`
}

// GetProduct godoc
//
//	@Summary		Get a product by ID
//	@Description	Returns a single product by its UUID.
//	@Tags			products
//	@Produce		json
//	@Param			id		path		string				true	"Product UUID"	Format(uuid)
//	@Success		200		{object}	GetProductResponse
//	@Failure		400		{object}	response.envelope	"Invalid UUID format"
//	@Failure		404		{object}	response.envelope	"Product not found"
//	@Failure		500		{object}	response.envelope	"Internal server error"
//	@Router			/products/{id} [get]
func NewGetProduct(u usecase.Executer[
	v1usecase.GetProductQuery,
	v1usecase.GetProductResult,
],
) handler.Route {
	return handler.Route{
		Method:  http.MethodGet,
		Path:    "/products/{id}",
		Handler: func(w http.ResponseWriter, r *http.Request) { getProduct(u, w, r) },
	}
}

func getProduct(u usecase.Executer[
	v1usecase.GetProductQuery,
	v1usecase.GetProductResult,
],
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		handler.Respond.BadRequest(w, r, "invalid id")
		return
	}

	qry := v1usecase.GetProductQuery{
		ID: id,
	}

	res, err := u.Execute(r.Context(), qry).ToGo()
	if err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	handler.Respond.Ok(w, r, GetProductResponse{
		ID:               res.Product.ID(),
		Name:             res.Product.Name(),
		Slug:             res.Product.Slug(),
		Description:      res.Product.Description(),
		ShortDescription: res.Product.ShortDescription(),
		DisplayImageURL:  res.Product.DisplayImageURL(),
		PriceCents:       res.Product.PriceCents(),
		Currency:         res.Product.Currency(),
		Status:           string(res.Product.Status()),
		CreatedAt:        res.Product.CreatedAt(),
		UpdatedAt:        res.Product.UpdatedAt(),
	})
}
