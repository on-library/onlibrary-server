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
		Username	string		`json:"username"`
		Password	string		`json:"password"`
	}

	RegisterRequest struct {
		LoginRequest
		Name		string		`json:"name"`
		Email		string		`json:"email"`
		Address		string		`json:"address"`
		City		string		`json:"city"`
		Province	string		`json:"province"`		
	}
)

func (controller AuthController) Routes() []common.Route {
	return []common.Route {
		{
			Method: echo.GET,
			Path: "/auth/profile",
			Handler: controller.Profile,
		},
		{
			Method:echo.POST,
			Path: "/auth/login",
			Handler: controller.Login,
		},
		{
			Method: echo.POST,
			Path: "/auth/register",
			Handler: controller.Register,
		},
	}
}

func (controller AuthController) Profile(c echo.Context) error {
	return c.String(http.StatusOK, "Profile")
}

func (controller AuthController) Login(c echo.Context) error {
	params := new(LoginRequest)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest,err)
	}

	
	return c.JSON(http.StatusOK, params)
}

func (controller AuthController) Register(c echo.Context) error {
	params := new(RegisterRequest)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadGateway, err)
	}


	return c.JSON(http.StatusOK, params)
}
