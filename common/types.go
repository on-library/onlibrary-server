package common

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	Route struct {
		Method		string
		Path		string
		Handler		echo.HandlerFunc
		Middleware	[]echo.MiddlewareFunc
	}

	Controller interface {
		Routes()	[]Route
	}

	CustomValidator struct {
		Validator	*validator.Validate
	}

	ValidationError struct {
		Namespace	string		`json:"namespace"`
		Field		string		`json:"field"`
		Tag			string		`json:"tag"`
		Message		string		`json:"message"`
	}

	ValidationErrors []ValidationError

	UserRole		string

	GeneralResponseJSON struct {
		Message 	string		`json:"message"`
	}
)

func (ve ValidationErrors) Error() string {
	sErrs := make([]string, len(ve))
	for i, validationError := range ve {
		sErrs[i] = validationError.Message
	}
	return strings.Join(sErrs, "\n")
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err:= cv.Validator.Struct(i); err!=nil{
		validationErrors := err.(validator.ValidationErrors)
		var errors ValidationErrors
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Namespace: e.Namespace(),
				Field: e.Field(),
				Tag: e.Tag(),
				Message: fmt.Sprintf("Field validation for '%s' failed on ther '%s' Tag", e.Field(), e.Tag()),
			})
		}
		return errors
	}

	return nil
}