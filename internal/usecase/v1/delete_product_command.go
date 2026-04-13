package v1usecase

import (
	"context"

	"github.com/google/uuid"

	result "github.com/selfshop-dev/lib-result"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/internal/usecase"
)

type DeleteProductCommand struct {
	ID uuid.UUID
}

type DeleteProductResult struct{}

type deleteProduct struct {
	r domain.ProductRepository
}

func NewDeleteProduct(r domain.ProductRepository) usecase.Executer[
	DeleteProductCommand, DeleteProductResult,
] {
	return &deleteProduct{r: r}
}

func (u *deleteProduct) Execute(ctx context.Context, cmd DeleteProductCommand) result.Value[DeleteProductResult] {
	return result.Map(
		result.Of(struct{}{}, u.r.Delete(ctx, cmd.ID)),
		func(struct{}) DeleteProductResult { return DeleteProductResult{} },
	)
}
