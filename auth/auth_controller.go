package auth

import (
	"fmt"
	"math/rand"
	"net/http"
	"onlibrary/auth/models"
	"onlibrary/common"
	"onlibrary/database"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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
		Name			string		`json:"name" validate:"required,min=6"`
		Email			string		`json:"email" validate:"required,email,min=6"`
		Nim				string		`json:"nim" validate:"required,min=6"`
		TanggalLahir	time.Time	`json:"tanggal_lahir" validate:"required"`
		Address			string		`json:"address" validate:"required"`	
	}

	VerifyRequest struct {
		Username	string		`json:"username" validate:"required"`
		Code		int			`json:"code" validate:"required"`
	}

	EditRequest struct {
		ID			uuid.UUID		`json:"id"`
		Name		string			`json:"name" validate:"required"`
		Email		string			`json:"email" validate:"required,email"`
		Address		string			`json:"address" validate:"required"`
	}


)

func (controller AuthController) Routes() []common.Route {
	return []common.Route {
		{
			Method: echo.GET,
			Path: "/auth/profile",
			Handler: controller.Profile,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleware()},
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
		{
			Method: echo.POST,
			Path: "/auth/verify",
			Handler: controller.VerifyAccount,
		},
		{
			Method: echo.POST,
			Path: "/auth/edit",
			Handler: controller.EditAccount,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleware()},
		},
		{
			Method: echo.GET,
			Path: "/auth/all",
			Handler: controller.GetAuths,
		},
	}
}

func (controller AuthController) Login(c echo.Context) error {
	db := database.GetInstance()
	params := new(LoginRequest)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest,err)
	}

	var user models.Auth

	var r = struct {
		common.GeneralResponseJSON
	}{
		GeneralResponseJSON: common.GeneralResponseJSON{Message:"Username/password tidak dikenali"},
	}

	if db.First(&user, "username = ?", params.Username);user.Username != params.Username{
		return c.JSON(http.StatusBadRequest, r)
	}

	if !CheckPasswordHash(params.Password, user.Password){
		return c.JSON(http.StatusBadRequest, r)
	}

	if user.IsVerify == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":"failed",
			"message":"Account is not valid. Please check your email for activation",
		})
	}

	claims := &common.JwtCustomClaims{
		ID: user.ID,
		Username: user.Username,
		Role: user.Role,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("SECRETKEY"))
	if err != nil {
		fmt.Print(err)
		return err
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"status":"success",
		"token":t,
	})
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
	newUser.Name = params.Name
	newUser.Role = 1
	newUser.Nim = params.Nim
	newUser.Email = params.Email
	newUser.Address = params.Address
	newUser.IsVerify = 0
	newUser.VerifyCode =rand.Intn(9999)


	var verifLink = "http://dev-onlibrary.s3-website-ap-southeast-2.amazonaws.com/#/auth/verify/" + strconv.Itoa(newUser.VerifyCode) + "?username="+newUser.Username

	var email = EmailInfo{
		Body: fmt.Sprintf("Harap kunjungi link <a href='%s'>ini</a> untuk verifikasi akun anda", verifLink),
		Subject: "Verifikasi akun",
		From: "onlibraryid@gmail.com",
		To: newUser.Email,
	}

	SendEmail(email)

	db.Create(&newUser)

	return c.JSON(http.StatusOK, params)
}



func (controller AuthController) Profile(c echo.Context) error {
	db:= database.GetInstance()

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)

	var userInfo models.Auth

	db.Preload("Rents.Book").Preload("Rents.Book.Genres").Preload("Rents.Book.Category").Preload("Rents.Book.Reviews").Preload("Rents",func(db *gorm.DB) *gorm.DB{
		return db.Order("rents.updated_at DESC")
	}).First(&userInfo,"username = ?",claims.Username)
	
	return c.JSON(http.StatusOK, echo.Map{
		"user":userInfo,
		
	})
}

func (controller AuthController) VerifyAccount (c echo.Context) error {
	db := database.GetInstance()
	params := new(VerifyRequest)

	if err:= c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, err)
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


	if err := db.First(&user,"username = ?", params.Username); err.Error !=nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message":"Username not found","status":"error"})
	}

	if user.VerifyCode != params.Code {
		return c.JSON(http.StatusBadRequest, echo.Map{"message":"Invalid code","status":"error"})
	}

	user.IsVerify = 1

	db.Save(&user)

	return c.JSON(http.StatusOK,echo.Map{
		"status":"success",
		"message":"Account verified!",
	})
}

func (controller AuthController) GetAuths(c echo.Context) error {
	db := database.GetInstance()
	var auths []models.Auth

	db.Find(&auths)
	

	var r = struct {
		common.GeneralResponseJSON
		Data []models.Auth		`json:"data"`
	}{}


	r.Message = "Success"
	r.Data = auths

	return c.JSON(http.StatusOK, r)
}

func (controller AuthController) EditAccount (c echo.Context) error {
	db:= database.GetInstance()
	params := new(EditRequest)

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	
	if err:=c.Bind(params);err!=nil{
		return c.JSON(http.StatusBadRequest,err)
		
	}

	var auth models.Auth

	if err := db.First(&auth, "id = ?", claims.ID); err.Error !=nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message":"User not found","status":"error"})
	}

	auth.Name = params.Name
	auth.Email = params.Email
	auth.Address = params.Address

	db.Save(&auth)

	return c.JSON(http.StatusOK,
		echo.Map{
			"message":"Account changed successfully",
			"status":"success",
			"data":auth,
		})
}