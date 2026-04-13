package v1usecase

import (
	"context"

	result "github.com/selfshop-dev/lib-result"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
)

type GetProductBySlugQuery struct {
	Slug string
}

type GetProductBySlugResult struct {
	Product *domain.Product
}

type getProductBySlug struct {
	r domain.ProductRepository
}

func NewGetProductBySlug(r domain.ProductRepository) usecase.Executer[
	GetProductBySlugQuery, GetProductBySlugResult,
] {
	return &getProductBySlug{r: r}
}

func (u *getProductBySlug) Execute(ctx context.Context, qry GetProductBySlugQuery) result.Value[GetProductBySlugResult] {
	return result.Map(
		result.Of(u.r.GetBySlug(ctx, qry.Slug)),
		func(p *domain.Product) GetProductBySlugResult { return GetProductBySlugResult{Product: p} },
	)
}
