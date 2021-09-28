package common

import (
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
		Namespace	string
		Field		string
		Tag			string
		Message		string
	}

	ValidationErrors []ValidationError

	UserRole		string
)