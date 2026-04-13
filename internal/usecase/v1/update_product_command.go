package v1usecase

import (
	"context"

	"github.com/google/uuid"

	result "github.com/selfshop-dev/lib-result"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
)

type UpdateProductCommand struct {
	Description      *string
	ShortDescription *string
	DisplayImageURL  *string
	Name             string
	Slug             string
	Currency         string
	PriceCents       int64
	ID               uuid.UUID
}

type UpdateProductResult struct {
	ID uuid.UUID
}

type updateProduct struct {
	r domain.ProductRepository
}

func NewUpdateProduct(r domain.ProductRepository) usecase.Executer[
	UpdateProductCommand, UpdateProductResult,
] {
	return &updateProduct{r: r}
}

func (u *updateProduct) Execute(ctx context.Context, cmd UpdateProductCommand) result.Value[UpdateProductResult] {
	data, err := domain.NewProduct(
		cmd.Name,
		cmd.Slug,
		cmd.Description,
		cmd.ShortDescription,
		cmd.DisplayImageURL,
		cmd.PriceCents,
		cmd.Currency,
		new(domain.ProductStatusActive),
	)
	if err != nil {
		return result.Err[UpdateProductResult](err)
	}
	return result.Map(
		result.Of(u.r.Update(ctx, cmd.ID, data)),
		func(p *domain.Product) UpdateProductResult { return UpdateProductResult{ID: p.ID()} },
	)
}
