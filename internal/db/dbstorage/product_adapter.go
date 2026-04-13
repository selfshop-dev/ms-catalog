package dbstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	apperr "github.com/selfshop-dev/lib-apperr"
	ctxval "github.com/selfshop-dev/lib-ctxval"

	"github.com/selfshop-dev/ms-catalog/internal/db/gen"
	"github.com/selfshop-dev/ms-catalog/internal/domain"
)

type productAdapter struct {
	g *gen.Queries
}

func NewProductAdapter(g *gen.Queries) domain.ProductRepository {
	return &productAdapter{g: g}
}

func (a *productAdapter) Create(ctx context.Context, data *domain.Product) (*domain.Product, error) {
	params := gen.CreateProductParams{
		Name:             data.Name(),
		Slug:             data.Slug(),
		Description:      data.Description(),
		ShortDescription: data.ShortDescription(),
		DisplayImageURL:  data.DisplayImageURL(),
		PriceCents:       data.PriceCents(),
		Currency:         data.Currency(),
		Status:           data.Status(),
	}
	row, err := a.q(ctx).CreateProduct(ctx, params)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, apperr.Conflictf("product with slug %s already exists", data.Slug())
		}
		return nil, err
	}
	return restoreProduct(&row), nil
}

func (a *productAdapter) Update(ctx context.Context, id uuid.UUID, data *domain.Product) (*domain.Product, error) {
	params := gen.UpdateProductParams{
		ID:               id,
		Name:             data.Name(),
		Slug:             data.Slug(),
		Description:      data.Description(),
		ShortDescription: data.ShortDescription(),
		DisplayImageURL:  data.DisplayImageURL(),
		PriceCents:       data.PriceCents(),
		Currency:         data.Currency(),
	}
	row, err := a.q(ctx).UpdateProduct(ctx, params)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, apperr.Conflictf("product with slug %s already exists", data.Slug())
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFoundf("product not found with ID %s", id)
		}
		return nil, err
	}
	return restoreProduct(&row), nil
}

func (a *productAdapter) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) (*domain.Product, error) {
	params := gen.UpdateStatusProductParams{
		ID:     id,
		Status: status,
	}
	row, err := a.q(ctx).UpdateStatusProduct(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFoundf("product not found with ID %s", id)
		}
		return nil, err
	}
	return restoreProduct(&row), nil
}

func (a *productAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	affected, err := a.q(ctx).DeleteProduct(ctx, id)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFoundf("product not found with ID %s", id)
	}
	return nil
}

func (a *productAdapter) GetListActiveProducts(ctx context.Context, limit, offset int32) ([]*domain.Product, error) {
	params := gen.GetListActiveProductsParams{
		Limit:  limit,
		Offset: offset,
	}
	rows, err := a.q(ctx).GetListActiveProducts(ctx, params)
	if err != nil {
		return nil, err
	}
	products := make([]*domain.Product, 0, len(rows))
	for i := range rows {
		products = append(products, restoreProduct(&rows[i]))
	}
	return products, nil
}

func (a *productAdapter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	row, err := a.q(ctx).GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFoundf("product not found with ID %s", id)
		}
		fmt.Printf("DEBUG GetByID error: %T %v\n", err, err)
		return nil, err
	}
	return restoreProduct(&row), nil
}

func (a *productAdapter) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	row, err := a.q(ctx).GetProductBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFoundf("product not found with slug %s", slug)
		}
		return nil, err
	}
	return restoreProduct(&row), nil
}

func (a *productAdapter) q(ctx context.Context) *gen.Queries { return ctxval.Or(ctx, a.g) }

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func restoreProduct(row *gen.Product) *domain.Product {
	return domain.RestoreProduct(
		row.ID,
		row.Name,
		row.Slug,
		row.Description,
		row.ShortDescription,
		row.DisplayImageURL,
		row.PriceCents,
		row.Currency,
		row.Status,
		row.CreatedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)
}
