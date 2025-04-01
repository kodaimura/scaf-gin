package main

import (
	"github.com/labstack/echo/v4"

	"goscaf/internal/router"
)

func main() {
	e := echo.New()
	router.Set(e)
	e.Logger.Fatal(e.Start(":3000"))
}