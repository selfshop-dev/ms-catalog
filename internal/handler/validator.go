package handler

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	validation "github.com/selfshop-dev/lib-validation"
)

var v *validator.Validate

func init() {
	v = validator.New(
		validator.WithRequiredStructEnabled(),
	)
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
}

func Validate(s any) error {
	if err := v.Struct(s); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
			return convertValidationErrors(ve)
		}
		return err
	}
	return nil
}

func convertValidationErrors(es validator.ValidationErrors) error {
	vc := validation.NewCollector("request")
	for _, e := range es {
		vc.Add(tagToFieldError(e))
	}
	return vc.Validation()
}

func tagToFieldError(e validator.FieldError) validation.FieldError {
	field := e.Field()
	param := e.Param()

	switch e.Tag() {
	case "required":
		return validation.Required(field)
	case "min":
		n, err := strconv.Atoi(param)
		if err != nil {
			return validation.Invalid(field, field+" is invalid")
		}
		return validation.TooShort(field, n)
	case "max":
		n, err := strconv.Atoi(param)
		if err != nil {
			return validation.Invalid(field, field+" is invalid")
		}
		return validation.TooLong(field, n)
	case "len":
		return validation.Invalid(field, fmt.Sprintf("%s must be exactly %s characters long", field, param))
	case "gte":
		return validation.OutOfRange(field, param, "∞")
	case "lte":
		return validation.OutOfRange(field, "0", param)
	case "url":
		return validation.Invalid(field, fmt.Sprintf("%s must be a valid URL", field))
	case "oneof":
		return validation.Invalid(field, fmt.Sprintf("%s must be one of: %s", field, param))
	default:
		return validation.Invalid(field, fmt.Sprintf("%s is invalid (%s)", field, e.Tag()))
	}
}
