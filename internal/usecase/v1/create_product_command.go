package v1usecase

import (
	"context"

	"github.com/google/uuid"

	result "github.com/selfshop-dev/lib-result"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
)

type CreateProductCommand struct {
	Description      *string
	ShortDescription *string
	DisplayImageURL  *string
	Name             string
	Slug             string
	Currency         string
	Status           domain.ProductStatus
	PriceCents       int64
}

type CreateProductResult struct {
	ID uuid.UUID
}

type createProduct struct {
	r domain.ProductRepository
}

func NewCreateProduct(r domain.ProductRepository) usecase.Executer[
	CreateProductCommand, CreateProductResult,
] {
	return &createProduct{r: r}
}

func (u *createProduct) Execute(ctx context.Context, cmd CreateProductCommand) result.Value[CreateProductResult] {
	data, err := domain.NewProduct(
		cmd.Name,
		cmd.Slug,
		cmd.Description,
		cmd.ShortDescription,
		cmd.DisplayImageURL,
		cmd.PriceCents,
		cmd.Currency,
		new(cmd.Status),
	)
	if err != nil {
		return result.Err[CreateProductResult](err)
	}
	return result.Map(
		result.Of(u.r.Create(ctx, data)),
		func(p *domain.Product) CreateProductResult { return CreateProductResult{ID: p.ID()} },
	)
}
