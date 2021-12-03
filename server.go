package main

import (
	"log"
	authModel "onlibrary/auth/models"
	bookModel "onlibrary/books/models"
	categoryModel "onlibrary/category/models"
	"onlibrary/common"
	"onlibrary/database"
	genreModel "onlibrary/genre/models"
	rentModel "onlibrary/rents/models"
	reviewModel "onlibrary/reviews/models"
	"onlibrary/routes"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main(){
	err := godotenv.Load(".env")
	if err!= nil {
		log.Fatal("Error loading .env file")
	}

	api := echo.New()
	api.Validator = &common.CustomValidator{Validator: validator.New()}

	api.Use(middleware.CORS())

	db := database.GetInstance()
	db.AutoMigrate(&bookModel.Book{})
	db.AutoMigrate(&authModel.Auth{})
	db.AutoMigrate(&reviewModel.Review{})
	db.AutoMigrate(&rentModel.Rent{})
	db.AutoMigrate(&genreModel.Genre{})
	db.AutoMigrate(&categoryModel.Category{})

	routes.DefineApiRoute(api)

	server := echo.New()
	server.Any("/*", func(c echo.Context) (err error) {
		req:= c.Request()
		res:= c.Response()
		if req.URL.Path[:4] == "/api" {
			api.ServeHTTP(res,req)
		}

		return
	})

	server.Logger.Fatal(server.Start(":8080"))


	
}