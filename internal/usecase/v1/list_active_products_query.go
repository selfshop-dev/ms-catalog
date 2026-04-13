package v1usecase

import (
	"context"

	result "github.com/selfshop-dev/lib-result"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
)

type GetActiveProductsQuery struct {
	Limit  int32
	Offset int32
}

type GetActiveProductsResult struct {
	Products []*domain.Product
}

type getActiveProducts struct {
	r domain.ProductRepository
}

func NewGetActiveProducts(r domain.ProductRepository) usecase.Executer[
	GetActiveProductsQuery, GetActiveProductsResult,
] {
	return &getActiveProducts{r: r}
}

func (u *getActiveProducts) Execute(ctx context.Context, qry GetActiveProductsQuery) result.Value[GetActiveProductsResult] {
	limit := qry.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := max(qry.Offset, 0)
	return result.Map(
		result.Of(u.r.GetListActiveProducts(ctx, limit, offset)),
		func(ps []*domain.Product) GetActiveProductsResult { return GetActiveProductsResult{Products: ps} },
	)
}
