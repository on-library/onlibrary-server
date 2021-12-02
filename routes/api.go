package routes

import (
	"onlibrary/auth"
	"onlibrary/books"
	"onlibrary/category"
	"onlibrary/common"
	"onlibrary/genre"
	"onlibrary/rents"
	"onlibrary/reviews"

	"github.com/labstack/echo/v4"
)

func DefineApiRoute(e *echo.Echo){
	controllers := []common.Controller{
		auth.AuthController{},
		books.BooksController{},
		reviews.ReviewController{},
		rents.RentController{},
		genre.GenreController{},
		category.CategoryController{},
	}
	var routes []common.Route
	for _, controller := range controllers {
		routes = append(routes, controller.Routes()...)
	}
	api := e.Group("/api")
	for _, route := range routes{
		switch route.Method {
		case echo.POST:
			{
				api.POST(route.Path, route.Handler, route.Middleware...)
				break
			}
		case echo.GET:
			{
				api.GET(route.Path, route.Handler, route.Middleware...)
				break
			}
		case echo.DELETE:
			{
				api.DELETE(route.Path, route.Handler, route.Middleware...)
				break
			}
		case echo.PUT:
			{
				api.PUT(route.Path, route.Handler, route.Middleware...)
				break
		 	}
		case echo.PATCH:
			{
				api.PATCH(route.Path, route.Handler, route.Middleware...)
				break
			}
		}
		
	}
}