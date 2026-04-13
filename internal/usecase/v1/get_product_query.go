package v1usecase

import (
	"context"

	"github.com/google/uuid"

	result "github.com/selfshop-dev/lib-result"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
)

type GetProductQuery struct {
	ID uuid.UUID
}

type GetProductResult struct {
	Product *domain.Product
}

type getProduct struct {
	r domain.ProductRepository
}

func NewGetProduct(r domain.ProductRepository) usecase.Executer[
	GetProductQuery, GetProductResult,
] {
	return &getProduct{r: r}
}

func (u *getProduct) Execute(ctx context.Context, qry GetProductQuery) result.Value[GetProductResult] {
	return result.Map(
		result.Of(u.r.GetByID(ctx, qry.ID)),
		func(p *domain.Product) GetProductResult { return GetProductResult{Product: p} },
	)
}
