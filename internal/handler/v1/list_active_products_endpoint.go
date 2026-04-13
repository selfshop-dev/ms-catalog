package v1handler

import (
	"net/http"
	"strconv"

	"github.com/selfshop-dev/ms-catalog/internal/handler"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
)

type GetActiveProductsResponse struct {
	Items  []GetProductResponse `json:"items"`
	Limit  int32                `json:"limit"`
	Offset int32                `json:"offset"`
}

// GetActiveProducts godoc
//
//	@Summary		List active products
//	@Description	Returns a paginated list of active products.
//	@Tags			products
//	@Produce		json
//	@Param			limit	query		int					false	"Max items to return (1–100, default 20)"
//	@Param			offset	query		int					false	"Number of items to skip (default 0)"
//	@Success		200		{object}	GetActiveProductsResponse
//	@Failure		500		{object}	response.envelope	"Internal server error"
//	@Router			/products/active [get]
func NewGetActiveProducts(u usecase.Executer[
	v1usecase.GetActiveProductsQuery,
	v1usecase.GetActiveProductsResult,
],
) handler.Route {
	return handler.Route{
		Method:  http.MethodGet,
		Path:    "/products/active",
		Handler: func(w http.ResponseWriter, r *http.Request) { getActiveProducts(u, w, r) },
	}
}

func getActiveProducts(u usecase.Executer[
	v1usecase.GetActiveProductsQuery,
	v1usecase.GetActiveProductsResult,
],
	w http.ResponseWriter,
	r *http.Request,
) {
	limit := int32(20)
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil &&
			n > 0 && n <= 100 {
			limit = int32(n)
		}
	}

	offset := int32(0)
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil &&
			n >= 0 {
			offset = int32(n)
		}
	}

	qry := v1usecase.GetActiveProductsQuery{
		Limit:  limit,
		Offset: offset,
	}

	res, err := u.Execute(r.Context(), qry).ToGo()
	if err != nil {
		handler.Respond.Error(w, r, err)
		return
	}

	items := make([]GetProductResponse, 0, len(res.Products))
	for _, p := range res.Products {
		items = append(items, GetProductResponse{
			ID:               p.ID(),
			Name:             p.Name(),
			Slug:             p.Slug(),
			Description:      p.Description(),
			ShortDescription: p.ShortDescription(),
			DisplayImageURL:  p.DisplayImageURL(),
			PriceCents:       p.PriceCents(),
			Currency:         p.Currency(),
			Status:           string(p.Status()),
			CreatedAt:        p.CreatedAt(),
			UpdatedAt:        p.UpdatedAt(),
		})
	}

	handler.Respond.Ok(w, r, GetActiveProductsResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	})
}
