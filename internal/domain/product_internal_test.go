package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	validation "github.com/selfshop-dev/lib-validation"
)

// Each validator is called directly to test its own error messages and boundary
// conditions in isolation, without going through NewProduct.
func newCollector() *validation.Collector {
	return validation.NewCollector("product")
}

func TestValidateProductName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid single char", input: "A", wantErr: false},
		{name: "valid 128 chars", input: strings.Repeat("x", 128), wantErr: false},
		{name: "empty string", input: "", wantErr: true},
		{name: "129 chars", input: strings.Repeat("x", 129), wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vc := newCollector()
			validateProductName(vc, tc.input)
			if tc.wantErr {
				require.Error(t, vc.Err(), "validateProductName(%q) expected error", tc.input)
			} else {
				require.NoError(t, vc.Err(), "validateProductName(%q) unexpected error", tc.input)
			}
		})
	}
}

func TestValidateProductSlug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		slug    string
		wantErr bool
	}{
		{name: "valid", slug: "my-product-1", wantErr: false},
		{name: "single letter", slug: "a", wantErr: false},
		{name: "single digit", slug: "1", wantErr: false},
		{name: "max length", slug: strings.Repeat("a", 128), wantErr: false},
		{name: "empty", slug: "", wantErr: true},
		{name: "too long", slug: strings.Repeat("a", 129), wantErr: true},
		{name: "starts with hyphen", slug: "-bad", wantErr: true},
		{name: "ends with hyphen", slug: "bad-", wantErr: true},
		{name: "uppercase", slug: "Bad", wantErr: true},
		{name: "underscore", slug: "bad_slug", wantErr: true},
		{name: "double hyphen", slug: "my--product", wantErr: false}, // regex allows it
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vc := newCollector()
			validateProductSlug(vc, tc.slug)
			if tc.wantErr {
				require.Error(t, vc.Err(), "validateProductSlug(%q) expected error", tc.slug)
			} else {
				require.NoError(t, vc.Err(), "validateProductSlug(%q) unexpected error", tc.slug)
			}
		})
	}
}

func TestValidateProductShortDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		input   *string
		wantErr bool
	}{
		{name: "nil", input: nil, wantErr: false},
		{name: "empty string", input: new(""), wantErr: false},
		{name: "exactly 256 chars", input: new(strings.Repeat("a", 256)), wantErr: false},
		{name: "257 chars", input: new(strings.Repeat("a", 257)), wantErr: true},
		{name: "256 chars after trimming whitespace", input: new("  " + strings.Repeat("a", 256) + "  "), wantErr: false},
		{name: "257 chars after trimming whitespace", input: new("  " + strings.Repeat("a", 257) + "  "), wantErr: true},
		{name: "only spaces under limit", input: new("  "), wantErr: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vc := newCollector()
			validateProductShortDescription(vc, tc.input)
			if tc.wantErr {
				require.Error(t, vc.Err())
			} else {
				require.NoError(t, vc.Err())
			}
		})
	}
}

func TestValidateProductDisplayImageURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		url     *string
		wantErr bool
	}{
		{name: "nil", url: nil, wantErr: false},
		{name: "https URL", url: new("https://cdn.example.com/a.jpg"), wantErr: false},
		{name: "http URL", url: new("http://cdn.example.com/a.jpg"), wantErr: false},
		{name: "ftp rejected", url: new("ftp://cdn.example.com/a.jpg"), wantErr: true},
		{name: "no scheme", url: new("cdn.example.com/a.jpg"), wantErr: true},
		{name: "empty string", url: new(""), wantErr: true},
		{name: "whitespace only", url: new("   "), wantErr: true},
		{name: "https after spaces", url: new("   https://ok.com/x.jpg"), wantErr: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vc := newCollector()
			validateProductDisplayImageURL(vc, tc.url)
			if tc.wantErr {
				require.Error(t, vc.Err(), "validateProductDisplayImageURL(%v) expected error", tc.url)
			} else {
				require.NoError(t, vc.Err(), "validateProductDisplayImageURL(%v) unexpected error", tc.url)
			}
		})
	}
}

func TestValidateProductPriceCents(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		priceCents int64
		wantErr    bool
	}{
		{name: "zero", priceCents: 0, wantErr: false},
		{name: "positive", priceCents: 1, wantErr: false},
		{name: "large", priceCents: 1_000_000_00, wantErr: false},
		{name: "minus one", priceCents: -1, wantErr: true},
		{name: "very negative", priceCents: -9999, wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vc := newCollector()
			validateProductPriceCents(vc, tc.priceCents)
			if tc.wantErr {
				require.Error(t, vc.Err(), "validateProductPriceCents(%d) expected error", tc.priceCents)
			} else {
				require.NoError(t, vc.Err(), "validateProductPriceCents(%d) unexpected error", tc.priceCents)
			}
		})
	}
}

func TestValidateProductCurrency(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		currency string
		wantErr  bool
	}{
		{name: "USD", currency: "USD", wantErr: false},
		{name: "EUR", currency: "EUR", wantErr: false},
		{name: "RUB", currency: "RUB", wantErr: false},
		{name: "lowercase", currency: "usd", wantErr: true},
		{name: "mixed case", currency: "Usd", wantErr: true},
		{name: "2 chars", currency: "US", wantErr: true},
		{name: "4 chars", currency: "USDD", wantErr: true},
		{name: "empty", currency: "", wantErr: true},
		{name: "with digit", currency: "U1D", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vc := newCollector()
			validateProductCurrency(vc, tc.currency)
			if tc.wantErr {
				require.Error(t, vc.Err(), "validateProductCurrency(%q) expected error", tc.currency)
			} else {
				require.NoError(t, vc.Err(), "validateProductCurrency(%q) unexpected error", tc.currency)
			}
		})
	}
}

func TestValidateProductStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		status  ProductStatus
		wantErr bool
	}{
		{name: "active", status: ProductStatusActive, wantErr: false},
		{name: "inactive", status: ProductStatusInactive, wantErr: false},
		{name: "draft", status: ProductStatusDraft, wantErr: false},
		{name: "archived", status: ProductStatusArchived, wantErr: false},
		{name: "unknown", status: ProductStatus("unknown"), wantErr: true},
		{name: "empty", status: ProductStatus(""), wantErr: true},
		{name: "mixed case active", status: ProductStatus("Active"), wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vc := newCollector()
			validateProductStatus(vc, tc.status)
			if tc.wantErr {
				require.Error(t, vc.Err(), "validateProductStatus(%q) expected error", tc.status)
			} else {
				require.NoError(t, vc.Err(), "validateProductStatus(%q) unexpected error", tc.status)
			}
		})
	}
}

func TestValidateProductName_ErrorFieldName(t *testing.T) {
	t.Parallel()

	vc := newCollector()
	validateProductName(vc, "")

	ve := vc.Validation()
	require.NotNil(t, ve)
	assert.NotEmpty(t, ve.Fields, "expected at least one field error")
}
