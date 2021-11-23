package books

import (
	"fmt"
	"net/http"
	"onlibrary/books/models"
	modelCategory "onlibrary/category/models"
	"onlibrary/common"
	"onlibrary/database"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type(
	BookController struct{

	}

	
	AddBookRequest struct {
		JudulBuku			string		`json:"judul_buku"`
		DeskripsiBuku		string		`json:"deskripsi_buku"`
		Penulis			string		`json:"penulis"`
		CategoryID		uuid.UUID		`json:"category_id"`
		Penerbit		string		`json:"penerbit"`
		TahunTerbit		time.Time	`json:"tahun_terbit"`
		Genres			[]string	`json:"genres"`	
		Stok			int		`json:"stok"`
		CreatedAt		time.Time	`json:"created_at"`
		UpdatedAt 		time.Time	`json:"updated_at"`
	}

	EditBookRequest struct {
		ID				uuid.UUID
		AddBookRequest
	}

	

)


func (controller BookController) Routes() []common.Route {
	return []common.Route {
		{
			Method: echo.GET,
			Path: "/book/:bookId",
			Handler: controller.GetBook,
		},
		{
			Method: echo.GET,
			Path: "/book/all",
			Handler: controller.GetBooks,
		},
		{
			Method: echo.GET,
			Path: "/book/filter",
			Handler: controller.FilterBooks,
		},
		{
			Method: echo.PUT,
			Path: "/book/edit",
			Handler: controller.EditBook,
		},
		{
			Method: echo.POST,
			Path: "/book/add",
			Handler: controller.AddBook,
		},
		{
			Method: echo.DELETE,
			Path: "/book/:bookId",
			Handler: controller.DeleteBook,
		},
	}
}

func (controller BookController) GetBooks(c echo.Context) error {
	db := database.GetInstance()
	var books []models.Book

	// TODO: Add filter using query params (title, author, publisher, category)
	db.Preload("Reviews").Find(&books)

	var r = struct {
		common.GeneralResponseJSON
		Data []models.Book		`json:"data"`
	}{}


	r.Message = "Success"
	r.Data = books

	return c.JSON(http.StatusOK, r)
}

func (controller BookController) GetBook(c echo.Context) error {
	bookId := c.Param("bookId")
	db := database.GetInstance()
	var book models.Book
	err := db.Preload("Reviews").First(&book, "id = ?",bookId).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}

	var r = struct {
		common.GeneralResponseJSON
		Data models.Book `json:"data"`
	}{
		GeneralResponseJSON:common.GeneralResponseJSON{Message: "Success"},
		Data: book,
	}
	return c.JSON(http.StatusOK, r)
}


func (controller BookController) AddBook(c echo.Context) error {
	params := new(AddBookRequest)

	if err := c.Bind(params); err!=nil{
		return c.JSON(http.StatusBadRequest, err)
	}

	var category modelCategory.Category 

	var newBook models.Book

	id := uuid.NewV1()


	newBook.ID = id;
	newBook.JudulBuku = params.JudulBuku
	newBook.TahunTerbit = params.TahunTerbit
	newBook.DeskripsiBuku = params.DeskripsiBuku
	newBook.BookCategoryID = params.CategoryID

	fmt.Println(params.Genres)

	db := database.GetInstance()

	if err:= db.First(&category, "category_id = ?", params.CategoryID); err.Error != nil {
		var r = struct {
			common.GeneralResponseJSON
		}{
			GeneralResponseJSON: common.GeneralResponseJSON{Message: "Category not found"},
		}
		fmt.Println(err.Error)
		return c.JSON(http.StatusBadRequest, r)
	}
	
	db.Model(&category).Association("Books").Append(&newBook)

	var r = struct {
		common.GeneralResponseJSON
		Data models.Book `json:"data"`
	}{
	GeneralResponseJSON:common.GeneralResponseJSON{Message: "Success"},
		Data: newBook,
	}
	return c.JSON(http.StatusOK, r)

}

func (controller BookController) EditBook(c echo.Context) error {
	db := database.GetInstance()
	params := new(EditBookRequest)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var book models.Book


	db.First(&book, "id = ?", params.ID)

	
	
	book.JudulBuku = params.JudulBuku
	book.DeskripsiBuku = params.DeskripsiBuku
	book.Penulis = params.Penulis
	book.BookCategoryID = params.CategoryID
	book.DeskripsiBuku = params.DeskripsiBuku
	book.Stok = params.Stok

	db.Save(&book)

	return c.JSON(http.StatusOK, book)
}

func (controller BookController) DeleteBook(c echo.Context) error {
	params := c.Param("bookId")
	db := database.GetInstance()

	var book models.Book

	db.Where("id = ?", params).Delete(&book)

	var r = struct {
		common.GeneralResponseJSON
		Id string
	}{
		GeneralResponseJSON:common.GeneralResponseJSON{Message: "Success"},
		Id: params,
	}
	return c.JSON(http.StatusOK, r)
}

func (controller BookController) FilterBooks(c echo.Context) error {
	//
	return c.JSON(http.StatusOK, "D")
}