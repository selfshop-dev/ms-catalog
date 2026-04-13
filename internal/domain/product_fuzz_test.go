package domain

// go test -fuzz=Fuzz -fuzztime=30s ./internal/domain/...

import (
	"strings"
	"testing"

	validation "github.com/selfshop-dev/lib-validation"
)

func FuzzValidateProductSlug(f *testing.F) {
	// seed corpus — known valid and invalid cases
	f.Add("widget")
	f.Add("")
	f.Add("-bad")
	f.Add("bad-")
	f.Add(strings.Repeat("a", 129))

	f.Fuzz(func(_ *testing.T, slug string) {
		vc := validation.NewCollector("product")
		// must never panic regardless of input
		validateProductSlug(vc, slug)
	})
}
