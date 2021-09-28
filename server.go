package main

import (
	"net/http"
	"onlibrary/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

func saveUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK,id)
}

func show(c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team : "+team+"member : "+member)
}

func main(){
	// e := echo.New()
	// e.Static("/static", "static")
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello World")
	// })

	// e.POST("/users", saveUser)
	// e.GET("/users/:id", getUser)
	// 	// e.PUT("/users/:id", updateUser)
	// 	// e.DELETE("/users/:id", deleteUser)
	
	// e.GET("/show", show)

	

	// e.Logger.Fatal(e.Start(":1323"))
	api := echo.New()

	api.Use(middleware.CORS())
	
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