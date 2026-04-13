package domain_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/selfshop-dev/ms-catalog/internal/domain"
	"github.com/selfshop-dev/ms-catalog/testhelpers"
)

func validProductArgs() (string, string, *string, *string, *string, int64, string, domain.ProductStatus) {
	return "My Product",
		"my-product",
		new("Full description"),
		new("Short desc"),
		new("https://example.com/image.jpg"),
		1000,
		"USD",
		domain.ProductStatusActive
}

func TestNewProduct_HappyPath(t *testing.T) {
	t.Parallel()

	name, slug, desc, short, img, price, cur, status := validProductArgs()
	p, err := domain.NewProduct(name, slug, desc, short, img, price, cur, new(status))

	require.NoError(t, err)
	assert.Equal(t, "My Product", p.Name())
	assert.Equal(t, "my-product", p.Slug())
	assert.Equal(t, desc, p.Description())
	assert.Equal(t, short, p.ShortDescription())
	assert.Equal(t, img, p.DisplayImageURL())
	assert.Equal(t, int64(1000), p.PriceCents())
	assert.Equal(t, "USD", p.Currency())
	assert.Equal(t, domain.ProductStatusActive, p.Status())
}

func TestNewProduct_NilOptionalFields(t *testing.T) {
	t.Parallel()

	p, err := domain.NewProduct("Name", "name", nil, nil, nil, 0, "EUR", new(domain.ProductStatusDraft))

	require.NoError(t, err)
	assert.Nil(t, p.Description())
	assert.Nil(t, p.ShortDescription())
	assert.Nil(t, p.DisplayImageURL())
}

func TestNewProduct_TrimsWhitespace(t *testing.T) {
	t.Parallel()

	p, err := domain.NewProduct("  Trimmed  ", "  my-slug  ", nil, nil, nil, 0, "USD", new(domain.ProductStatusDraft))

	require.NoError(t, err)
	assert.Equal(t, "Trimmed", p.Name())
	assert.Equal(t, "my-slug", p.Slug())
}

func TestNewProduct_ZeroPriceCents(t *testing.T) {
	t.Parallel()

	_, err := domain.NewProduct("Name", "name", nil, nil, nil, 0, "USD", new(domain.ProductStatusActive))
	require.NoError(t, err)
}

func TestNewProduct_NameValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid min length", input: "A", wantErr: false},
		{name: "valid max length", input: strings.Repeat("a", 128), wantErr: false},
		{name: "empty name", input: "", wantErr: true},
		{name: "name too long", input: strings.Repeat("a", 129), wantErr: true},
		{name: "whitespace only trims to empty", input: "   ", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewProduct(tc.input, "valid-slug", nil, nil, nil, 0, "USD", new(domain.ProductStatusDraft))
			if tc.wantErr {
				require.Error(t, err, "NewProduct(%q) expected error", tc.input)
			} else {
				require.NoError(t, err, "NewProduct(%q) unexpected error", tc.input)
			}
		})
	}
}

func TestNewProduct_SlugValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slug    string
		wantErr bool
	}{
		{name: "valid lowercase", slug: "my-product", wantErr: false},
		{name: "valid alphanumeric", slug: "product123", wantErr: false},
		{name: "valid single char", slug: "a", wantErr: false},
		{name: "valid max length", slug: strings.Repeat("a", 128), wantErr: false},
		{name: "empty slug", slug: "", wantErr: true},
		{name: "too long", slug: strings.Repeat("a", 129), wantErr: true},
		{name: "starts with hyphen", slug: "-invalid", wantErr: true},
		{name: "ends with hyphen", slug: "invalid-", wantErr: true},
		{name: "uppercase letters", slug: "Invalid", wantErr: true},
		{name: "spaces", slug: "my product", wantErr: true},
		{name: "special chars", slug: "my_product", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewProduct("Name", tc.slug, nil, nil, nil, 0, "USD", new(domain.ProductStatusDraft))
			if tc.wantErr {
				require.Error(t, err, "NewProduct(slug=%q) expected error", tc.slug)
			} else {
				require.NoError(t, err, "NewProduct(slug=%q) unexpected error", tc.slug)
			}
		})
	}
}

func TestNewProduct_ShortDescriptionValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		input   *string
		wantErr bool
	}{
		{name: "nil allowed", input: nil, wantErr: false},
		{name: "empty string allowed", input: new(""), wantErr: false},
		{name: "exactly max length", input: new(strings.Repeat("a", 256)), wantErr: false},
		{name: "exceeds max length", input: new(strings.Repeat("a", 257)), wantErr: true},
		{name: "whitespace within limit", input: new("  short  "), wantErr: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewProduct("Name", "name", nil, tc.input, nil, 0, "USD", new(domain.ProductStatusDraft))
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewProduct_DisplayImageURLValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		url     *string
		wantErr bool
	}{
		{name: "nil allowed", url: nil, wantErr: false},
		{name: "valid https", url: new("https://example.com/img.png"), wantErr: false},
		{name: "valid http", url: new("http://example.com/img.jpg"), wantErr: false},
		{name: "missing scheme", url: new("example.com/img.png"), wantErr: true},
		{name: "ftp scheme rejected", url: new("ftp://example.com/img.png"), wantErr: true},
		{name: "empty string rejected", url: new(""), wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewProduct("Name", "name", nil, nil, tc.url, 0, "USD", new(domain.ProductStatusDraft))
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewProduct_PriceCentsValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		priceCents int64
		wantErr    bool
	}{
		{name: "zero", priceCents: 0, wantErr: false},
		{name: "positive", priceCents: 999, wantErr: false},
		{name: "large value", priceCents: 1_000_000_00, wantErr: false},
		{name: "negative", priceCents: -1, wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewProduct("Name", "name", nil, nil, nil, tc.priceCents, "USD", new(domain.ProductStatusDraft))
			if tc.wantErr {
				require.Error(t, err, "NewProduct(priceCents=%d) expected error", tc.priceCents)
			} else {
				require.NoError(t, err, "NewProduct(priceCents=%d) unexpected error", tc.priceCents)
			}
		})
	}
}

func TestNewProduct_CurrencyValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		currency string
		wantErr  bool
	}{
		{name: "valid USD", currency: "USD", wantErr: false},
		{name: "valid EUR", currency: "EUR", wantErr: false},
		{name: "lowercase rejected", currency: "usd", wantErr: true},
		{name: "too short", currency: "US", wantErr: true},
		{name: "too long", currency: "USDD", wantErr: true},
		{name: "empty", currency: "", wantErr: true},
		{name: "with digit", currency: "US1", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewProduct("Name", "name", nil, nil, nil, 0, tc.currency, new(domain.ProductStatusDraft))
			if tc.wantErr {
				require.Error(t, err, "NewProduct(currency=%q) expected error", tc.currency)
			} else {
				require.NoError(t, err, "NewProduct(currency=%q) unexpected error", tc.currency)
			}
		})
	}
}

func TestNewProduct_StatusValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		status  domain.ProductStatus
		wantErr bool
	}{
		{name: "active", status: domain.ProductStatusActive, wantErr: false},
		{name: "inactive", status: domain.ProductStatusInactive, wantErr: false},
		{name: "draft", status: domain.ProductStatusDraft, wantErr: false},
		{name: "archived", status: domain.ProductStatusArchived, wantErr: false},
		{name: "unknown status", status: domain.ProductStatus("unknown"), wantErr: true},
		{name: "empty status", status: domain.ProductStatus(""), wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewProduct("Name", "name", nil, nil, nil, 0, "USD", new(tc.status))
			if tc.wantErr {
				require.Error(t, err, "NewProduct(status=%q) expected error", tc.status)
			} else {
				require.NoError(t, err, "NewProduct(status=%q) unexpected error", tc.status)
			}
		})
	}
}

func TestNewProduct_MultipleValidationErrors(t *testing.T) {
	t.Parallel()

	// empty name + bad currency + bad status — all three errors must surface
	_, err := domain.NewProduct("", "valid-slug", nil, nil, nil, 0, "bad", new(domain.ProductStatus("nope")))
	require.Error(t, err)
}

func TestRestoreProduct_AccessorsReturnStoredValues(t *testing.T) {
	t.Parallel()

	id := testhelpers.UUIDNew()

	createdAt := testhelpers.TimeNow()
	updatedAt := testhelpers.TimeNow()

	p := domain.RestoreProduct(
		id, "Restored", "restored-slug",
		new("desc"), new("short"),
		new("https://cdn.example.com/img.jpg"),
		500, "GBP",
		domain.ProductStatusInactive,
		createdAt,
		updatedAt,
		new(updatedAt),
	)

	assert.Equal(t, id, p.ID())
	assert.Equal(t, "Restored", p.Name())
	assert.Equal(t, "restored-slug", p.Slug())
	assert.Equal(t, "GBP", p.Currency())
	assert.Equal(t, int64(500), p.PriceCents())
	assert.Equal(t, domain.ProductStatusInactive, p.Status())
	assert.Equal(t, createdAt, p.CreatedAt())
	assert.Equal(t, updatedAt, p.UpdatedAt())
	assert.NotNil(t, p.DeletedAt())
}
