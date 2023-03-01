package controller

import (
	"github.com/labstack/echo/v4"
)

// add middleware basic plus customizable
func Start(e *echo.Echo) {
	//starts an http server

	//e.Use(middleware.BasicAuth())
	e.GET("/user", GetUser)
	e.GET("/users", getAllUser)
	e.POST("/create", CreateUser)
	e.PUT("/user", UpdateUser)
}
