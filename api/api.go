package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//InitAPI starts web api for maestro
func InitAPI() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	group := e.Group("/v1.0.0")
	// Routes
	group.GET("/reprocessing/pending", getPendingReprocessing)
	group.POST("/reprocess/top", reprocessTop)
	group.POST("/reprocess/top/skip", reprocessTop)
	// Start server
	e.Logger.Fatal(e.Start(":8089"))
}