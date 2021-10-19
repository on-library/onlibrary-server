package auth

import (
	"net/http"
	"onlibrary/auth/models"
	"onlibrary/common"
	"onlibrary/database"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type(
	AuthController struct{

	}

	LoginRequest struct {
		Username	string		`json:"username" validate:"required"`
		Password	string		`json:"password" validate:"required"`
	}

	RegisterRequest struct {
		LoginRequest
		Name		string		`json:"name" validate:"required"`
		Email		string		`json:"email" validate:"required,email"`
		Address		string		`json:"address" validate:"required"`
		City		string		`json:"city" validate:"required"`
		Province	string		`json:"province" validate:"required"`		
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
	db := database.GetInstance()
	params := new(RegisterRequest)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadGateway, err)
	}

	if err:= c.Validate(params); err != nil {
		var r = struct {
			common.GeneralResponseJSON
			Errors error		`json:"errors"`
		}{
			GeneralResponseJSON: common.GeneralResponseJSON{Message: "Error"},
			Errors: err,
		}
		return echo.NewHTTPError(http.StatusBadRequest, r)
	}

	var user models.Auth

	if db.First(&user, "username = ?", params.Username);user.Username == params.Username  {
		var r = struct {
			common.GeneralResponseJSON
		}{GeneralResponseJSON: common.GeneralResponseJSON{Message: "Username already exist!"}}
		return c.JSON(http.StatusBadRequest, r)
	}


	hashedPassword, _ := HashPassword(params.Password)
	newID := uuid.NewV1()

	var newUser models.Auth

	newUser.ID = newID
	newUser.Username = params.Username
	newUser.Password = hashedPassword
	newUser.Role = "1"
	newUser.Email = params.Email
	newUser.City = params.City
	newUser.Province = params.Province
	newUser.Address = params.Address

	db.Create(&newUser)

	return c.JSON(http.StatusOK, params)
}
