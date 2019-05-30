package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func heartBeatHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func main() {
	e := echo.New()

	// Monitoring handlers
	e.GET("/ping", heartBeatHandler)

	e.Logger.Fatal(e.Start(":1323"))
}