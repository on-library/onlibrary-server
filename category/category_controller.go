package category

import (
	"net/http"
	"onlibrary/category/models"
	"onlibrary/common"
	"onlibrary/database"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type(
	CategoryController struct{

	}

	CategoryAddRequest struct {
		Nama		string	`json:"nama"`
	}

)


func (controller CategoryController) Routes() []common.Route{
	return []common.Route{
		{
			Method: echo.GET,
			Path: "/category",
			Handler: controller.GetCategories,
		},
		{
			Method: echo.POST,
			Path: "/category/add",
			Handler: controller.AddCategory,
		},
		
	}
}

func (controller CategoryController) GetCategories (c echo.Context) error {
	db := database.GetInstance()
	var categories []models.Category

	// db.Preload("Books").Find(&categories)
	if err := db.Find(&categories); err.Error!=nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":"error",
			"message":"Message not found",
		})
	}


	var r = struct {
		common.GeneralResponseJSON
		Data []models.Category		`json:"data"`
	}{}


	r.Message = "Success"
	r.Data = categories

	return c.JSON(http.StatusOK, r)

}

func (controller CategoryController) AddCategory (c echo.Context) error {
	db:= database.GetInstance()


	params := new(CategoryAddRequest)

	if err := c.Bind(params); err!=nil{
		return c.JSON(http.StatusBadRequest, err)
	}

	var newCategory models.Category

	id := uuid.NewV1()


	newCategory.CategoryID = id;
	newCategory.Nama = params.Nama


	db.Create(&newCategory)

	var r = struct {
		common.GeneralResponseJSON
		Data models.Category `json:"data"`
	}{
	GeneralResponseJSON:common.GeneralResponseJSON{Message: "Success"},
		Data: newCategory,
	}
	return c.JSON(http.StatusOK, r)
}