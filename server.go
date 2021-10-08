package main

import (
	"log"
	"onlibrary/books/models"
	"onlibrary/database"
	"onlibrary/routes"

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

	api.Use(middleware.CORS())

	db := database.GetInstance()
	db.AutoMigrate(&models.Book{})

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

	server.Logger.Fatal(server.Start(":1323"))


	
}