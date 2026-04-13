package domain

// go test -bench=. -benchmem -benchtime=3s ./internal/domain/...

import (
	"strings"
	"testing"

	validation "github.com/selfshop-dev/lib-validation"
)

func BenchmarkNewProduct(b *testing.B) {
	desc := "some description"
	b.ResetTimer()
	for b.Loop() {
		_, _ = NewProduct(
			"Widget", "widget",
			&desc, nil, nil,
			1000, "USD",
			new(ProductStatusDraft),
		)
	}
}

func BenchmarkValidateProductSlug(b *testing.B) {
	slugs := []string{
		"valid-slug-123",
		strings.Repeat("a", 128),
		"-invalid",
		"",
	}
	b.ResetTimer()
	for b.Loop() {
		vc := validation.NewCollector("product")
		validateProductSlug(vc, slugs[b.N%len(slugs)])
	}
}
