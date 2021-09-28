package auth

import (
	"net/http"
	"onlibrary/common"

	"github.com/labstack/echo/v4"
)

type(
	AuthController struct{

	}

	LoginRequest struct {

	}
)

func (controller AuthController) Routes() []common.Route {
	return []common.Route {
		{
			Method: echo.GET,
			Path: "/auth/profile",
			Handler: controller.Profile,
		},
	}
}

func (controller AuthController) Profile(c echo.Context) error {
	return c.String(http.StatusOK, "Profile")
}

