package main

import (
	"github.com/labstack/echo"
	"go_webservice/handlers"
	"net/http"
)

func heartBeatHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func main() {
	e := echo.New()

	// Monitoring handlers
	e.GET("/ping", heartBeatHandler)

	// Customised handlers
	e.POST("/file", handlers.CreateNewFileHandler)
	e.GET("/file", handlers.GetFileContentHandler)
	e.PUT("/file", handlers.ReplaceFileContentHandler)
	e.DELETE("/file", handlers.RemoveFileHandler)

	e.GET("/folder", handlers.GetFolderStatsHandler)

	e.Logger.Fatal(e.Start(":1323"))
}