package handler

import (
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	validation "github.com/selfshop-dev/lib-validation"
)

// stubFieldError is a test double for validator.FieldError.
// It allows constructing arbitrary tag/field/param combinations that are
// impossible to produce through the public Validate function — for example,
// a non-numeric param for "min" or "max" tags.
type stubFieldError struct {
	tag   string
	field string
	param string
}

func (s stubFieldError) Tag() string                      { return s.tag }
func (s stubFieldError) ActualTag() string                { return s.tag }
func (s stubFieldError) Field() string                    { return s.field }
func (s stubFieldError) Param() string                    { return s.param }
func (s stubFieldError) Namespace() string                { return "" }
func (s stubFieldError) StructNamespace() string          { return "" }
func (s stubFieldError) StructField() string              { return "" }
func (s stubFieldError) Value() any                       { return nil }
func (s stubFieldError) Kind() reflect.Kind               { return reflect.Invalid }
func (s stubFieldError) Type() reflect.Type               { return nil }
func (s stubFieldError) Translate(_ ut.Translator) string { return "" }
func (s stubFieldError) Error() string                    { return "" }

// compile-time check: stubFieldError satisfies validator.FieldError.
var _ validator.FieldError = stubFieldError{}

func TestTagToFieldError_min_invalid_param(t *testing.T) {
	t.Parallel()

	// go-playground/validator always produces numeric params for "min",
	// so this branch is unreachable through Validate — tested directly here.
	fe := tagToFieldError(stubFieldError{
		tag:   "min",
		field: "some_field",
		param: "not-a-number",
	})

	assert.Equal(t, validation.CodeInvalid, fe.Code)
	assert.Contains(t, fe.Message, "some_field is invalid")
}

func TestTagToFieldError_max_invalid_param(t *testing.T) {
	t.Parallel()

	// go-playground/validator always produces numeric params for "max",
	// so this branch is unreachable through Validate — tested directly here.
	fe := tagToFieldError(stubFieldError{
		tag:   "max",
		field: "some_field",
		param: "not-a-number",
	})

	assert.Equal(t, validation.CodeInvalid, fe.Code)
	assert.Contains(t, fe.Message, "some_field is invalid")
}
