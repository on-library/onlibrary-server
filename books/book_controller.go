package books

import (
	"net/http"
	"onlibrary/books/models"
	"onlibrary/common"

	"github.com/labstack/echo/v4"
)

type(
	BookController struct{

	}

	AddBookRequest struct {
		Title		string	`json:"Title"`
		Description	string 	`json:"Description"`
	}

)

var listBook = []models.Book{{Title: "A",Description: "2"},{Title: "d",Description: "2"}}

func (controller BookController) Routes() []common.Route {
	return []common.Route {
		{
			Method: echo.GET,
			Path: "/book/all",
			Handler: controller.GetBooks,
		},
		{
			Method: echo.POST,
			Path: "/book/add",
			Handler: controller.AddBook,
		},
	}
}

func (controller BookController) GetBooks(c echo.Context) error {
	return c.JSON(http.StatusOK, listBook)
}

func (controller BookController) AddBook(c echo.Context) error {
	params := new(AddBookRequest)

	if err := c.Bind(params); err!=nil{
		return c.JSON(http.StatusBadRequest, err)
	}

	var newBook models.Book

	newBook.Title = params.Title
	newBook.Description = params.Description

	var s = append(listBook, newBook)

	return c.JSON(http.StatusOK, s)
}