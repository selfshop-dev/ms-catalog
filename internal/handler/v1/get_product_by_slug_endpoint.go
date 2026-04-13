package v1handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/selfshop-dev/ms-catalog/internal/handler"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
)

// GetProductBySlug godoc
//
//	@Summary		Get a product by slug
//	@Description	Returns a single product by its slug.
//	@Tags			products
//	@Produce		json
//	@Param			slug	path		string				true	"Product slug"
//	@Success		200		{object}	GetProductResponse
//	@Failure		400		{object}	response.envelope	"Invalid slug"
//	@Failure		404		{object}	response.envelope	"Product not found"
//	@Failure		500		{object}	response.envelope	"Internal server error"
//	@Router			/products/{slug}/by-slug [get]
func NewGetProductBySlug(u usecase.Executer[
	v1usecase.GetProductBySlugQuery,
	v1usecase.GetProductBySlugResult,
],
) handler.Route {
	return handler.Route{
		Method:  http.MethodGet,
		Path:    "/products/{slug}/by-slug",
		Handler: func(w http.ResponseWriter, r *http.Request) { getProductBySlug(u, w, r) },
	}
}

func getProductBySlug(u usecase.Executer[
	v1usecase.GetProductBySlugQuery,
	v1usecase.GetProductBySlugResult,
],
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := chi.URLParam(r, "slug")
	// chi guarantees a non-empty slug when the route matches,
	// but guard defensively against misconfigured routers
	if slug == "" {
		handler.Respond.BadRequest(w, r, "invalid slug")
		return
	}

	qry := v1usecase.GetProductBySlugQuery{
		Slug: slug,
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
