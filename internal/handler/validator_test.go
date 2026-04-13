package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	validation "github.com/selfshop-dev/lib-validation"

	"github.com/selfshop-dev/ms-catalog/internal/handler"
)

// testStruct covers all tag branches exercised through the public Validate function.
type testStruct struct {
	Required string `json:"required_field" validate:"required"`
	Min      string `json:"min_field"      validate:"min=3"`
	Max      string `json:"max_field"      validate:"max=5"`
	Len      string `json:"len_field"      validate:"len=3"`
	Gte      int    `json:"gte_field"      validate:"gte=10"`
	Lte      int    `json:"lte_field"      validate:"lte=100"`
	URL      string `json:"url_field"      validate:"omitempty,url"`
	OneOf    string `json:"oneof_field"    validate:"omitempty,oneof=foo bar"`
	Email    string `json:"email_field"    validate:"omitempty,email"` // triggers default branch
}

func TestValidate_valid(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Required: "ok",
		Min:      "abc",
		Max:      "ab",
		Len:      "abc",
		Gte:      10,
		Lte:      50,
	}

	err := handler.Validate(s)
	require.NoError(t, err)
}

func TestValidate_required(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Min: "abc",
		Max: "ab",
		Len: "abc",
		Gte: 10,
		Lte: 50,
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("required_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeRequired, fe.Code)
}

func TestValidate_min(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Required: "ok",
		Min:      "ab", // shorter than min=3
		Max:      "ab",
		Len:      "abc",
		Gte:      10,
		Lte:      50,
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("min_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeTooShort, fe.Code)
	assert.Equal(t, 3, fe.Meta["min"])
}

func TestValidate_max(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Required: "ok",
		Min:      "abc",
		Max:      "toolong", // longer than max=5
		Len:      "abc",
		Gte:      10,
		Lte:      50,
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("max_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeTooLong, fe.Code)
	assert.Equal(t, 5, fe.Meta["max"])
}

func TestValidate_len(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Required: "ok",
		Min:      "abc",
		Max:      "ab",
		Len:      "ab", // not exactly len=3
		Gte:      10,
		Lte:      50,
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("len_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeInvalid, fe.Code)
	assert.Contains(t, fe.Message, "3")
}

func TestValidate_gte(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Required: "ok",
		Min:      "abc",
		Max:      "ab",
		Len:      "abc",
		Gte:      5, // less than gte=10
		Lte:      50,
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("gte_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeOutOfRange, fe.Code)
}

func TestValidate_lte(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Required: "ok",
		Min:      "abc",
		Max:      "ab",
		Len:      "abc",
		Gte:      10,
		Lte:      200, // greater than lte=100
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("lte_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeOutOfRange, fe.Code)
}

func TestValidate_url(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Required: "ok",
		Min:      "abc",
		Max:      "ab",
		Len:      "abc",
		Gte:      10,
		Lte:      50,
		URL:      "not-a-url",
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("url_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeInvalid, fe.Code)
	assert.Contains(t, fe.Message, "valid URL")
}

func TestValidate_oneof(t *testing.T) {
	t.Parallel()

	s := testStruct{
		Required: "ok",
		Min:      "abc",
		Max:      "ab",
		Len:      "abc",
		Gte:      10,
		Lte:      50,
		OneOf:    "baz", // not in oneof=foo bar
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("oneof_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeInvalid, fe.Code)
	assert.Contains(t, fe.Message, "foo bar")
}

func TestValidate_default_branch(t *testing.T) {
	t.Parallel()

	// email tag is not handled explicitly — falls through to default branch
	s := testStruct{
		Required: "ok",
		Min:      "abc",
		Max:      "ab",
		Len:      "abc",
		Gte:      10,
		Lte:      50,
		Email:    "not-an-email",
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	fe, ok := ve.First("email_field")
	require.True(t, ok)
	assert.Equal(t, validation.CodeInvalid, fe.Code)
	assert.Contains(t, fe.Message, "email")
}

func TestValidate_no_json_tag_uses_field_name(t *testing.T) {
	t.Parallel()

	type req struct {
		ProductName string `validate:"required"`
	}

	err := handler.Validate(req{})
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	_, ok = ve.First("ProductName")
	assert.True(t, ok, "expected field name to be 'ProductName' when json tag is absent")
}

func TestValidate_json_tag_dash_uses_field_name(t *testing.T) {
	t.Parallel()

	type req struct {
		ProductName string `json:"-" validate:"required"`
	}

	err := handler.Validate(req{})
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	_, ok = ve.First("ProductName")
	assert.True(t, ok, "expected field name to be 'ProductName' when json tag is '-'")
}

func TestValidate_json_field_name_used(t *testing.T) {
	t.Parallel()

	type req struct {
		ProductName string `json:"product_name" validate:"required"`
	}

	err := handler.Validate(req{})
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)

	_, ok = ve.First("product_name")
	assert.True(t, ok, "expected field name to be 'product_name' from json tag, not 'ProductName'")
}

func TestValidate_non_struct_returns_raw_error(t *testing.T) {
	t.Parallel()

	// passing nil triggers validator.InvalidValidationError — not ValidationErrors,
	// so the raw error is returned unwrapped
	err := handler.Validate(nil)
	require.Error(t, err)

	_, ok := validation.As(err)
	assert.False(t, ok, "expected raw error, not validation.Error")
}

func TestValidate_multiple_errors(t *testing.T) {
	t.Parallel()

	s := testStruct{
		// Required is empty
		Min: "a",       // too short
		Max: "toolong", // too long
		Len: "ab",      // wrong length
		Gte: 1,         // below gte
		Lte: 200,       // above lte
	}

	err := handler.Validate(s)
	require.Error(t, err)

	ve, ok := validation.As(err)
	require.True(t, ok)
	assert.GreaterOrEqual(t, len(ve.Fields), 2)
}
