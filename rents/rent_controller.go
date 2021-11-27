package rents

import (
	"fmt"
	"net/http"
	modelAuth "onlibrary/auth/models"
	modelBook "onlibrary/books/models"
	"onlibrary/common"
	"onlibrary/database"
	"onlibrary/rents/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// status peminjaman = 1 => buku menunggu diambil
// status peminjaman = 2 => buku sedang dipinjam
// status peminjaman = 3 => buku permintaan
// status peminjaman = 4=> buku selesai
//

type(
	RentController struct{

	}

	RentAddRequest struct {
		BookID			uuid.UUID		`json:"book_id"`
		CreatedAt		time.Time	`json:"created_at"`
		UpdatedAt 		time.Time	`json:"updated_at"`
		DeskirpsiPeminjaman string	`json:"deskripsi_peminjaman"`
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

	var visitor modelAuth.Auth
	var newRent models.Rent
	var book modelBook.Book

	if err:= db.First(&book, "book_id = ?", params.BookID);err.Error!=nil{
		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"Book not found",
		})
	}

	if err:= db.First(&visitor, "id = ?", claims.ID); err.Error != nil {
		var r = struct {
			common.GeneralResponseJSON
		}{
			GeneralResponseJSON: common.GeneralResponseJSON{Message: "User not found"},
		}
		fmt.Println(err.Error)
		return c.JSON(http.StatusBadRequest, r)
	}
	
	id := uuid.NewV1()
	newRent.PinjamID = id
	newRent.BookRentID = params.BookID
	newRent.UserRef = claims.ID
	newRent.TanggalPinjam = time.Now()
	newRent.TanggalPengembalian = time.Now().AddDate(0,0,7)
	newRent.StatusPinjam = 0
	newRent.IsExtendConfirm = 0
	newRent.Denda = 500
	newRent.DeskripsiPeminjaman = params.DeskirpsiPeminjaman
	newRent.IsExtendConfirm = 0
	newRent.AlasanPerpanjangan = ""

	
	db.Model(&visitor).Association("Rents").Append(&newRent)

	db.Preload("Reviews").Preload("Category").Preload("Genres").First(&book, "book_id = ?",params.BookID)

	fmt.Println(&book)

	d:=db.Model(&book).Association("Rents").Append(&newRent)
	
	fmt.Println(d)

	return c.JSON(http.StatusOK,echo.Map{"message":"Rent added","data":newRent})
}

func (controller RentController) ConfirmRentBook(c echo.Context) error {
	db := database.GetInstance()

	params := new(ConfirmRentRequest)

	if err := c.Bind(params);err!=nil{
		return c.JSON(http.StatusBadRequest, err)
	}

	var rent models.Rent

	result := db.First(&rent,"pinjam_id = ?",params.RentID)
	if(result.Error != nil){
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Pinjam ID not found"})
	}
	
	rent.StatusPinjam = 1

	db.Save(&rent)

	return c.JSON(http.StatusOK,echo.Map{"message":"Rent confirmed", "rent_id":rent.PinjamID,"rent_status":rent.PinjamID})
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

	if(rent.StatusPinjam != 0){
		return c.JSON(http.StatusBadRequest, echo.Map{"message":"Failed"})
	}

	rent.StatusPinjam = -1
	
	return c.JSON(http.StatusOK, echo.Map{"message":"Rent declined","rent_id":params.RentID, "rent_status":rent.PinjamID})
}

// func (controller RentController) RentReturnRequest(c echo.Context) error 