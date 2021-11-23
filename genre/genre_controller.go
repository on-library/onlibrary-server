package genre

import (
	"net/http"
	"onlibrary/common"
	"onlibrary/database"
	"onlibrary/genre/models"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type(
	GenreController struct{

	}

	GenreAddRequest struct {
		Nama		string	`json:"nama"`
	}

)


func (controller GenreController) Routes() []common.Route{
	return []common.Route{
		{
			Method: echo.GET,
			Path: "/genre",
			Handler: controller.GetGenres,
		},
		{
			Method: echo.POST,
			Path: "/genre/add",
			Handler: controller.AddGenre,
		},
		
	}
}

func (controller GenreController) GetGenres (c echo.Context) error {
	db := database.GetInstance()
	var genres []models.Genre

	db.Find(&genres)

	var r = struct {
		common.GeneralResponseJSON
		Data []models.Genre		`json:"data"`
	}{}


	r.Message = "Success"
	r.Data = genres

	return c.JSON(http.StatusOK, r)

}

func (controller GenreController) AddGenre (c echo.Context) error {
	db:= database.GetInstance()


	params := new(GenreAddRequest)

	if err := c.Bind(params); err!=nil{
		return c.JSON(http.StatusBadRequest, err)
	}

	var newGenre models.Genre

	id := uuid.NewV1()


	newGenre.GenreID = id;
	newGenre.Nama = params.Nama


	db.Create(&newGenre)

	var r = struct {
		common.GeneralResponseJSON
		Data models.Genre `json:"data"`
	}{
	GeneralResponseJSON:common.GeneralResponseJSON{Message: "Success"},
		Data: newGenre,
	}
	return c.JSON(http.StatusOK, r)
}

// func (controller RentController) RentBook(c echo.Context) error {
// 	db := database.GetInstance()

// 	user := c.Get("user").(*jwt.Token)
// 	// fmt.Print(user)
// 	claims := user.Claims.(*common.JwtCustomClaims)
// 	// fmt.Print(claims)

// 	params:= new(RentAddRequest)

// 	if err := c.Bind(params);err!=nil{
// 		return c.JSON(http.StatusBadRequest,err)
// 	}

// 	var newRent models.Rent
	
// 	id := uuid.NewV1()
// 	newRent.ID = id
// 	newRent.BookID = params.BookID
// 	newRent.UserID = claims.ID
// 	newRent.RentAt = time.Now()
// 	newRent.EndAt = time.Now().AddDate(0,0,7)
// 	newRent.RentStatus = 0
// 	newRent.IsExtendConfirm = 0
	
// 	db.Create(&newRent)

// 	return c.JSON(http.StatusOK,echo.Map{"message":"Rent added","data":newRent})
// }

// func (controller RentController) ConfirmRentBook(c echo.Context) error {
// 	db := database.GetInstance()

// 	params := new(ConfirmRentRequest)

// 	if err := c.Bind(params);err!=nil{
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	var rent models.Rent

// 	result := db.First(&rent,"id = ?",params.RentID)
// 	if(result.Error != nil){
// 		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Rent ID not found"})
// 	}
	
// 	rent.RentStatus = 1

// 	db.Save(&rent)

// 	return c.JSON(http.StatusOK,echo.Map{"message":"Rent confirmed", "rent_id":rent.ID,"rent_status":rent.RentStatus})
// }

// func (controller RentController) DeclineRentBook(c echo.Context) error {
// 	db:= database.GetInstance()

// 	params := new(DeclineRentRequest)

// 	if err := c.Bind(params);err!=nil{
// 		return c.JSON(http.StatusBadRequest,err)
// 	}

// 	var rent models.Rent

// 	if result := db.First(&rent,"id = ?",params.RentID); result.Error != nil {
// 		return echo.NewHTTPError(http.StatusNotFound,echo.Map{"message":"Rent ID not found"})
// 	}

// 	if(rent.RentStatus != 0){
// 		return c.JSON(http.StatusBadRequest, echo.Map{"message":"Failed"})
// 	}

// 	rent.RentStatus = -1
	
// 	return c.JSON(http.StatusOK, echo.Map{"message":"Rent declined","rent_id":rent.ID, "rent_status":rent.RentStatus})
// }

// // func (controller RentController) RentReturnRequest(c echo.Context) error 