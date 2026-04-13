package v1usecase

import (
	"context"

	"github.com/google/uuid"

	result "github.com/selfshop-dev/lib-result"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
)

type UpdateStatusProductCommand struct {
	Status domain.ProductStatus
	ID     uuid.UUID
}

type UpdateStatusProductResult struct {
	ID uuid.UUID
}

type updateStatusProduct struct {
	r domain.ProductRepository
}

func NewUpdateStatusProduct(r domain.ProductRepository) usecase.Executer[
	UpdateStatusProductCommand, UpdateStatusProductResult,
] {
	return &updateStatusProduct{r: r}
}

func (u *updateStatusProduct) Execute(ctx context.Context, cmd UpdateStatusProductCommand) result.Value[UpdateStatusProductResult] {
	return result.Map(
		result.Of(u.r.UpdateStatus(ctx, cmd.ID, cmd.Status)),
		func(p *domain.Product) UpdateStatusProductResult { return UpdateStatusProductResult{ID: p.ID()} },
	)
}
