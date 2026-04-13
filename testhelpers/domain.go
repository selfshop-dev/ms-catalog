package testhelpers

import (
	"time"

	"github.com/google/uuid"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
)

// RestoredProduct returns a minimal [domain.Product] with the given id for use in mock returns.
func RestoredProduct(id uuid.UUID) *domain.Product {
	return domain.RestoreProduct(
		id,
		"Widget", "widget",
		nil, nil, nil,
		0, "USD",
		domain.ProductStatusDraft,
		time.Now(), time.Now(), nil,
	)
}
