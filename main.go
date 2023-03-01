package main

import (
	"github.com/labstack/echo/v4"

	"pro/controller"
)

// write comments where ever applicable
// customize according to your need
func main() {
	controller.DB()
	var e = echo.New()
	controller.Start(e)
	e.Logger.Fatal(e.Start(":8000"))
}
