package books

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"onlibrary/books/models"
	modelCategory "onlibrary/category/models"
	"onlibrary/common"
	"onlibrary/database"
	modelGenre "onlibrary/genre/models"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type(
	BooksController struct{

	}

	
	AddBookRequest struct {
		JudulBuku		string		`json:"judul_buku"`
		DeskripsiBuku	string		`json:"deskripsi_buku"`
		Penulis			string		`json:"penulis"`
		CategoryID		uuid.UUID	`json:"category_id"`
		Penerbit		string		`json:"penerbit"`
		TahunTerbit		time.Time	`json:"tahun_terbit"`
		Genres			[]string	`json:"genres"`	
		Stok			int			`json:"stok"`
		CreatedAt		time.Time	`json:"created_at"`
		UpdatedAt 		time.Time	`json:"updated_at"`
		ImageBase		[]byte		`json:"image_base"`
	}

	FindBookRequest struct {
		BookTitle		string		`json:"book_title"`
	}

	EditBookRequest struct {
		ID				uuid.UUID
		AddBookRequest
	}

	

)


func (controller BooksController) Routes() []common.Route {
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
			Method: echo.POST,
			Path: "/book/find",
			Handler: controller.FindBook,
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

func (controller BooksController) GetBooks(c echo.Context) error {
	db := database.GetInstance()
	var books []models.Book

	// TODO: Add filter using query params (title, author, publisher, category)
	db.Preload("Reviews").Preload("Category").Preload("Genres").Find(&books)
	

	var r = struct {
		common.GeneralResponseJSON
		Data []models.Book		`json:"data"`
	}{}


	r.Message = "Success"
	r.Data = books

	return c.JSON(http.StatusOK, r)
}

func (controller BooksController) FindBook(c echo.Context) error {
	db := database.GetInstance()	
	var books []models.Book

	params := new(FindBookRequest)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Preload("Reviews").Preload("Category").Preload("Genres").Where("judul_buku LIKE ?","%"+params.BookTitle+"%").Find(&books); err.Error !=nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message":"book not found",
			"status":"error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":"success",
		"data":books,
	})
}

func (controller BooksController) GetBook(c echo.Context) error {
	bookId := c.Param("bookId")
	db := database.GetInstance()
	var book models.Book
	err := db.Preload("Reviews").Preload("Category").Preload("Genres").First(&book, "book_id = ?",bookId).Error
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


func (controller BooksController) AddBook(c echo.Context) error {
	params := new(AddBookRequest)

	if err := c.Bind(params); err!=nil{
		return c.JSON(http.StatusBadRequest, err)
	}

	var category modelCategory.Category 

	var newBook models.Book

	id := uuid.NewV1()

	jpgI,err := png.Decode(bytes.NewReader(params.ImageBase))
	if err!=nil {
		fmt.Println("ini errorn",err)
		return err;
	}


	urlSaved := "static/" + params.JudulBuku +".png"

	dst, err := os.Create(urlSaved)
	if err!=nil{
		fmt.Println(err)
		return err;
	}

	if err := jpeg.Encode(dst,jpgI,nil);err!=nil{
		fmt.Println("set")
		return err;
	}
	dst.Close()


	newBook.BookId = id;
	newBook.JudulBuku = params.JudulBuku
	newBook.DeskripsiBuku = params.DeskripsiBuku
	newBook.BookCategoryID = params.CategoryID
	newBook.Penerbit = params.Penerbit
	newBook.Penulis = params.Penulis
	newBook.TahunTerbit = params.TahunTerbit
	newBook.Stok = params.Stok
	newBook.StokAwal = params.Stok
	newBook.CreatedAt = time.Now()
	newBook.UpdatedAt = time.Now()
	newBook.ImgUrl = urlSaved


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
	db.Create(newBook)
	
	db.Model(&category).Association("Books").Append(&newBook)
	
	for i :=0;i<len(params.Genres);i++{
		var genre modelGenre.Genre
		fmt.Println(params.Genres[i])
		genre = modelGenre.Genre{GenreID: uuid.NewV4(),Nama: params.Genres[i],GenreBookID: newBook.BookId}
		db.Create(genre)
	}

	var r = struct {
		common.GeneralResponseJSON
		Data models.Book `json:"data"`
	}{
	GeneralResponseJSON:common.GeneralResponseJSON{Message: "Success"},
		Data: newBook,
	}
	return c.JSON(http.StatusOK, r)

}

func (controller BooksController) EditBook(c echo.Context) error {
	db := database.GetInstance()
	params := new(EditBookRequest)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var book models.Book


	db.First(&book, "book_id = ?", params.ID)
	
	book.JudulBuku = params.JudulBuku
	book.DeskripsiBuku = params.DeskripsiBuku
	book.Penulis = params.Penulis
	book.BookCategoryID = params.CategoryID
	book.DeskripsiBuku = params.DeskripsiBuku
	book.Stok = params.Stok

	db.Save(&book)

	return c.JSON(http.StatusOK, book)
}

func (controller BooksController) DeleteBook(c echo.Context) error {
	params := c.Param("bookId")
	db := database.GetInstance()

	var book models.Book

	if err:= db.First(&book, "book_id = ?", params); err.Error != nil {
		var r = struct {
			common.GeneralResponseJSON
		}{
			GeneralResponseJSON: common.GeneralResponseJSON{Message: "Book not found"},
		}
		fmt.Println(err.Error)
		return c.JSON(http.StatusBadRequest, r)
	}

	// for i :=0;i<len(book.Genres);i++{
	// 	var genre modelGenre.Genre
	// 	db.Where("genre_id = ?", book.Genres[i].GenreID).Delete(&genre)
	// }

	db.Model(&book).Association("Genres").Clear()

	

	db.Where("book_id = ?", params).Delete(&book)

	// fmt.Println(d.Error)

	var r = struct {
		common.GeneralResponseJSON
		Id string
	}{
		GeneralResponseJSON:common.GeneralResponseJSON{Message: "Success"},
		Id: params,
	}
	return c.JSON(http.StatusOK, r)
}

func (controller BooksController) FilterBooks(c echo.Context) error {
	//
	return c.JSON(http.StatusOK, "D")
}