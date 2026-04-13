package domain

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	validation "github.com/selfshop-dev/lib-validation"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusArchived ProductStatus = "archived"
)

const (
	productNameMinLen = 1
	productNameMaxLen = 128

	productSlugMinLen = 1
	productSlugMaxLen = 128

	productShortDescriptionMaxLen = 256
)

var (
	productSlugRe     = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
	productImageURLRe = regexp.MustCompile(`^https?://`)
	productCurrencyRe = regexp.MustCompile(`^[A-Z]{3}$`)
)

type ProductRepository interface {
	Create(ctx context.Context, data *Product) (*Product, error)
	Update(ctx context.Context, id uuid.UUID, data *Product) (*Product, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status ProductStatus) (*Product, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetListActiveProducts(ctx context.Context, limit, offset int32) ([]*Product, error)

	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetBySlug(ctx context.Context, slug string) (*Product, error)
}

//go:generate mockgen -typed -destination=../mocks/mock_product_repository.go -package=mocks github.com/selfshop-dev/ms-catalog/internal/domain ProductRepository

type Product struct {
	id               uuid.UUID
	name             string
	slug             string
	description      *string
	shortDescription *string
	displayImageURL  *string
	priceCents       int64
	currency         string
	status           ProductStatus
	createdAt        time.Time
	updatedAt        time.Time
	deletedAt        *time.Time
}

func NewProduct(
	name string,
	slug string,
	description *string,
	shortDescription *string,
	displayImageURL *string,
	priceCents int64,
	currency string,
	// status is optional: when nil the field is not validated and the zero
	// value is left unset — use this for update-style construction where
	// the caller does not intend to change the product's status.
	status *ProductStatus,
) (*Product, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)

	vc := validation.NewCollector("product")

	validateProductName(vc, name)
	validateProductSlug(vc, slug)
	validateProductShortDescription(vc, shortDescription)
	validateProductDisplayImageURL(vc, displayImageURL)
	validateProductPriceCents(vc, priceCents)
	validateProductCurrency(vc, currency)

	if status != nil {
		validateProductStatus(vc, *status)
	}

	if ve := vc.Validation(); ve != nil {
		return nil, ve
	}

	p := &Product{
		name:             name,
		slug:             slug,
		description:      description,
		shortDescription: shortDescription,
		displayImageURL:  displayImageURL,
		priceCents:       priceCents,
		currency:         currency,
	}
	if status != nil {
		p.status = *status
	}

	return p, nil
}

func (p *Product) ID() uuid.UUID             { return p.id }
func (p *Product) Name() string              { return p.name }
func (p *Product) Slug() string              { return p.slug }
func (p *Product) Description() *string      { return p.description }
func (p *Product) ShortDescription() *string { return p.shortDescription }
func (p *Product) DisplayImageURL() *string  { return p.displayImageURL }
func (p *Product) PriceCents() int64         { return p.priceCents }
func (p *Product) Currency() string          { return p.currency }
func (p *Product) Status() ProductStatus     { return p.status }
func (p *Product) CreatedAt() time.Time      { return p.createdAt }
func (p *Product) UpdatedAt() time.Time      { return p.updatedAt }
func (p *Product) DeletedAt() *time.Time     { return p.deletedAt }

func RestoreProduct(
	id uuid.UUID,
	name string,
	slug string,
	description *string,
	shortDescription *string,
	displayImageURL *string,
	priceCents int64,
	currency string,
	status ProductStatus,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *Product {
	return &Product{
		id:               id,
		name:             name,
		slug:             slug,
		description:      description,
		shortDescription: shortDescription,
		displayImageURL:  displayImageURL,
		priceCents:       priceCents,
		currency:         currency,
		status:           status,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
		deletedAt:        deletedAt,
	}
}

func validateProductName(vc *validation.Collector, name string) {
	if len(name) < productNameMinLen {
		vc.Add(validation.TooShort("name", productNameMinLen))
	}
	if len(name) > productNameMaxLen {
		vc.Add(validation.TooLong("name", productNameMaxLen))
	}
}

func validateProductSlug(vc *validation.Collector, slug string) {
	if len(slug) < productSlugMinLen {
		vc.Add(validation.TooShort("slug", productSlugMinLen))
	}
	if len(slug) > productSlugMaxLen {
		vc.Add(validation.TooLong("slug", productSlugMaxLen))
	}
	if !productSlugRe.MatchString(slug) {
		vc.Add(validation.Invalid("slug", "product slug must contain only lowercase letters, digits, and hyphens"))
	}
}

func validateProductShortDescription(vc *validation.Collector, shortDescription *string) {
	if shortDescription == nil {
		return
	}
	if trimmed := strings.TrimSpace(*shortDescription); len(trimmed) > productShortDescriptionMaxLen {
		vc.Add(validation.TooLong("short_description", productShortDescriptionMaxLen))
	}
}

func validateProductDisplayImageURL(vc *validation.Collector, displayImageURL *string) {
	if displayImageURL == nil {
		return
	}
	if trimmed := strings.TrimSpace(*displayImageURL); !productImageURLRe.MatchString(trimmed) {
		vc.Add(validation.Invalid("display_image_url", "must be a valid URL starting with http:// or https://"))
	}
}

func validateProductPriceCents(vc *validation.Collector, priceCents int64) {
	if priceCents < 0 {
		vc.Add(validation.OutOfRange("price_cents", 0, "∞"))
	}
}

func validateProductCurrency(vc *validation.Collector, currency string) {
	if !productCurrencyRe.MatchString(currency) {
		vc.Add(validation.Invalid("currency", "must be a valid 3-letter ISO currency code"))
	}
}

func validateProductStatus(vc *validation.Collector, s ProductStatus) {
	switch s {
	case ProductStatusActive, ProductStatusInactive, ProductStatusDraft, ProductStatusArchived:
	default:
		vc.Add(validation.Invalid("status", "must be one of: active, inactive, draft, archived"))
	}
}
