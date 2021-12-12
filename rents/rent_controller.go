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
	"gorm.io/gorm/clause"
)

// status peminjaman = 1 => buku menunggu diambil
// status peminjaman = 2 => buku sedang dipinjam
// status peminjaman = 3 => buku permintaan
// status peminjaman = 4=> buku selesai
//

type(
	RentController struct{

	}

	InfoUserRentAll struct{
		models.Rent
		username	string
		name		string
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

	TakeRentRequest struct {
		RentID		string		`json:"rent_id"`
	}

	ExtendRentRequest struct {
		RentID				string		`json:"rent_id"`
		AlasanPerpanjangan	string		`json:"alasan_perpanjangan"`

	}

	ConfirmExtendRentRequest struct {
		RentID			string		`json:"rent_id"`
	}

	DeclineExtendRentRequest struct {
		RentID			string		`json:"rent_id"`
		AlasanPenolakanPepanjangan 	string `json:"alasan_penolakan_perpanjangan"`
	}

	DeclineRentRequest struct {
		RentID			string		`json:"rent_id"`
	}

)


func (controller RentController) Routes() []common.Route{
	return []common.Route{
		{
			Method: echo.GET,
			Path: "/rent/all",
			Handler: controller.GetRents,
		},
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
			Method:echo.POST,
			Path: "/rent/take",
			Handler: controller.TakeRent,
		},
		{
			Method:echo.POST,
			Path:"/rent/extend",
			Handler: controller.ExtendRent,
		},
		{
			Method: echo.POST,
			Path: "/rent/extend/confirm",
			Handler: controller.ConfirmExtendRent,
		},
		{
			Method: echo.POST,
			Path: "/rent/extend/decline",
			Handler: controller.DeclineExtendRent,
		},
		{
			Method: echo.POST,
			Path: "/rent/decline",
			Handler: controller.DeclineRentBook,	
		},
		{
			Method: echo.POST,
			Path: "/rent/return",
			Handler: controller.ReturnRent,
		},
	}
}

func (controller RentController) GetRents(c echo.Context) error {
	db:= database.GetInstance()

	type RentWithName struct {
		models.Rent
		User modelAuth.Auth		`json:"user"`
	}
	
	var rents []models.Rent
	var rentWithName []RentWithName
	var user modelAuth.Auth

	

	db.
			Preload("Book.Category").
			Preload(clause.Associations).
			Order("tanggal_pinjam desc").
			Find(&rents)

	for i:=0;i<len(rents);i++{
		
		if err := db.Select("username").First(&user, "id = ?", rents[i].AuthID);err.Error!=nil{
			return c.JSON(http.StatusBadRequest, echo.Map{"message":"user id not found","status":"error"})
		}

		rentWithName = append(rentWithName, RentWithName{User: user, Rent: rents[i]})
	}


	return c.JSON(http.StatusOK, echo.Map{
		"message":"success",
		"data":rentWithName,
	})
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

	if book.Stok <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message":"Stok buku habis",
			"status":"failed",
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
	newRent.TanggalPinjam = time.Now()
	newRent.AuthID = claims.ID
	newRent.TanggalPengembalian = time.Now().AddDate(0,0,7)
	newRent.StatusPinjam = 0
	newRent.IsExtendConfirm = 0
	newRent.Denda = 0
	newRent.DeskripsiPeminjaman = params.DeskirpsiPeminjaman
	newRent.IsExtendConfirm = 0
	newRent.AlasanPerpanjangan = ""

	

	db.Model(&visitor).Association("Rents").Append(&newRent)

	db.Preload("Reviews").Preload("Category").Preload("Genres").First(&book, "book_id = ?",params.BookID)


	db.Model(&book).Association("Rents").Append(&newRent)

	

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

	return c.JSON(http.StatusOK,echo.Map{"message":"Rent confirmed", "rent_id":rent.PinjamID,"rent_status":rent.StatusPinjam})
}

func (controller RentController) TakeRent(c echo.Context) error {
	db:= database.GetInstance()

	params := new(TakeRentRequest)

	if err:= c.Bind(params);err!=nil{
		return c.JSON(http.StatusBadRequest, err)
	}

	var rent models.Rent

	fmt.Println(params.RentID)

	if err:= db.Preload(clause.Associations).First(&rent, "pinjam_id = ?", params.RentID);err.Error!= nil {
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Pinjam ID not found"})
	}
	
	rent.StatusPinjam = 2
	rent.Book.Stok = rent.Book.Stok - 1
	
	db.Save(&rent)

	return c.JSON(http.StatusOK, echo.Map{"message":"Book picked up by user", "rent_id":rent.PinjamID,"rent_status":rent.StatusPinjam,"stok":rent.Book.Stok})

}

func (controller RentController) ExtendRent(c echo.Context) error {
	db:= database.GetInstance()

	params:= new(ExtendRentRequest)

	if err:= c.Bind(params);err!=nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var rent models.Rent

	fmt.Println(params.RentID)

	if err:= db.Preload(clause.Associations).First(&rent, "pinjam_id = ?", params.RentID);err.Error!= nil {
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Pinjam ID not found"})
	}

	if rent.IsExtendConfirm== 1{
		return c.JSON(http.StatusBadRequest,echo.Map{"message":"Sudah pernah melakukan perpanjangan", "status":"failed"})
	}

	rent.StatusPinjam = 3
	rent.AlasanPerpanjangan = params.AlasanPerpanjangan
	db.Save(&rent)

	return c.JSON(http.StatusOK, echo.Map{"message":"Konfirmasi perpanjangan telah dikirim","rent_id":rent.PinjamID,"rent_status":rent.StatusPinjam})
}

func (controller RentController) ConfirmExtendRent(c echo.Context) error {
	db := database.GetInstance()

	params := new(ConfirmExtendRentRequest)

	if err := c.Bind(params); err!=nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var rent models.Rent

	if err := db.Preload(clause.Associations).First(&rent, "pinjam_id = ?", params.RentID); err.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Pinjam ID not found"})
	}

	rent.StatusPinjam = 2
	rent.IsExtendConfirm = 1
	rent.TanggalPengembalian = rent.TanggalPengembalian.AddDate(0,0,7)

	db.Save(&rent)

	return c.JSON(http.StatusOK, echo.Map{"message":"Perpanjangan di konfirmasi"})
}

func (controller RentController) DeclineExtendRent(c echo.Context) error {
	db := database.GetInstance()

	params := new(DeclineExtendRentRequest)

	if err := c.Bind(params);err!= nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var rent models.Rent

	if err := db.First(&rent, "pinjam_id = ?", params.RentID);err.Error!=nil{
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Pinjam ID not found"})
	}

	if time.Now().Before(rent.TanggalPengembalian) {
		rent.StatusPinjam = 2
		rent.IsExtendConfirm = 0
		rent.AlasanPenolakanPepanjangan = params.AlasanPenolakanPepanjangan
	} 

	db.Save(&rent)


	return c.JSON(http.StatusOK,echo.Map{"message":"Rent declined","rent_id":params.RentID, "rent_status":rent.StatusPinjam})
}

func (controller RentController) DeclineRentBook(c echo.Context) error {
	db:= database.GetInstance()

	params := new(DeclineRentRequest)

	if err := c.Bind(params);err!=nil{
		return c.JSON(http.StatusBadRequest,err)
	}

	var rent models.Rent

	if result := db.First(&rent,"pinjam_id = ?",params.RentID); result.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Rent ID not found"})
	}

	if(rent.StatusPinjam != 0){
		return c.JSON(http.StatusBadRequest, echo.Map{"message":"Failed"})
	}

	
	rent.StatusPinjam = -3
	
	db.Save(&rent)
	return c.JSON(http.StatusOK, echo.Map{"message":"Rent declined","rent_id":params.RentID, "rent_status":rent.StatusPinjam})
}

func (controller RentController) ReturnRent(c echo.Context) error  {
	db:= database.GetInstance()

	params := new(RentReturnRequest)

	if err := c.Bind(params); err!=nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var rent models.Rent

	if err := db.Preload(clause.Associations).First(&rent, "pinjam_id = ?", params.RentID); err.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Pinjam ID not found"})
	}

	dateNow := time.Now()

	if dateNow.Before(rent.TanggalPengembalian) {
		rent.StatusPinjam = -3
		rent.Denda = 0
	} else {
		rent.StatusPinjam = -3
		rent.Denda = 500 * (rent.TanggalPengembalian.Day() - dateNow.Day())
	}
	
	rent.TanggalPengembalianFinish = &dateNow
	rent.Book.Stok += 1
	db.Save(&rent)

	return c.JSON(http.StatusOK, echo.Map{"message":"Buku berhasil dikembalikan","rent_id":rent.PinjamID,"book":rent.Book})
	
}