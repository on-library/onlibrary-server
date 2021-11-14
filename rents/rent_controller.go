package rents

import (
	"net/http"
	"onlibrary/common"
	"onlibrary/database"
	"onlibrary/rents/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type(
	RentController struct{

	}

	RentAddRequest struct {
		BookID			uuid.UUID		`json:"book_id"`
		CreatedAt		time.Time	`json:"created_at"`
		UpdatedAt 		time.Time	`json:"updated_at"`
	}

	RentInformation struct {
		RentID		string			`json:"rent_id"`
	}
	RentReturnRequest struct {
		RentID		string 			`json:"rent_id"`
	}

	ConfirmRentRequest struct {
		RentID			string		`json:"rent_id"`
	}

	DeclineRentRequest struct {
		RentID			string		`json:"rent_id"`
	}

)


func (controller RentController) Routes() []common.Route{
	return []common.Route{
		{
			Method: echo.POST,
			Path: "/rent/add",
			Handler:controller.RentBook,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleware()},
		},
		{
			Method: echo.POST,
			Path: "/rent/confirm",
			Handler: controller.ConfirmRentBook,
		},
		{
			Method: echo.POST,
			Path: "/rent/decline",
			Handler: controller.DeclineRentBook,	
		},
	}
}



func (controller RentController) RentBook(c echo.Context) error {
	db := database.GetInstance()

	user := c.Get("user").(*jwt.Token)
	// fmt.Print(user)
	claims := user.Claims.(*common.JwtCustomClaims)
	// fmt.Print(claims)

	params:= new(RentAddRequest)

	if err := c.Bind(params);err!=nil{
		return c.JSON(http.StatusBadRequest,err)
	}

	var newRent models.Rent
	
	id := uuid.NewV1()
	newRent.ID = id
	newRent.BookID = params.BookID
	newRent.UserID = claims.ID
	newRent.RentAt = time.Now()
	newRent.EndAt = time.Now().AddDate(0,0,7)
	newRent.RentStatus = 0
	newRent.IsExtendConfirm = 0
	
	db.Create(&newRent)

	return c.JSON(http.StatusOK,echo.Map{"message":"Rent added","data":newRent})
}

func (controller RentController) ConfirmRentBook(c echo.Context) error {
	db := database.GetInstance()

	params := new(ConfirmRentRequest)

	if err := c.Bind(params);err!=nil{
		return c.JSON(http.StatusBadRequest, err)
	}

	var rent models.Rent

	result := db.First(&rent,"id = ?",params.RentID)
	if(result.Error != nil){
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Rent ID not found"})
	}
	
	rent.RentStatus = 1

	db.Save(&rent)

	return c.JSON(http.StatusOK,echo.Map{"message":"Rent confirmed", "rent_id":rent.ID,"rent_status":rent.RentStatus})
}

func (controller RentController) DeclineRentBook(c echo.Context) error {
	db:= database.GetInstance()

	params := new(DeclineRentRequest)

	if err := c.Bind(params);err!=nil{
		return c.JSON(http.StatusBadRequest,err)
	}

	var rent models.Rent

	if result := db.First(&rent,"id = ?",params.RentID); result.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Rent ID not found"})
	}

	if(rent.RentStatus != 0){
		return c.JSON(http.StatusBadRequest, echo.Map{"message":"Failed"})
	}

	rent.RentStatus = -1
	
	return c.JSON(http.StatusOK, echo.Map{"message":"Rent declined","rent_id":rent.ID, "rent_status":rent.RentStatus})
}

// func (controller RentController) RentReturnRequest(c echo.Context) error 