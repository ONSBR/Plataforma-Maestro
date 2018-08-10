package api

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type H map[string]interface{}

//InitAPI starts web api for maestro
func InitAPI() {
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = errorHandler
	group := e.Group("/v1.0.0")
	// Routes
	group.GET("/reprocessing/:systemId/pending", getPendingReprocessing)
	group.GET("/reprocessing/:systemId/find", queryReprocessing)
	group.POST("/reprocessing/:id/approve", approveReprocessing)
	group.POST("/reprocessing/:id/skip", reprocessingSkip)
	group.GET("/gateway/:systemId/proceed", eventGatekeeper)
	group.POST("/handler/persist", startSystemPersistHandler)
	group.POST("/persist/sync", syncPersist)
	group.POST("/reprocessing/failure", reprocessingFailure)
	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func errorHandler(err error, c echo.Context) {
	if errJ := c.JSON(400, map[string]string{"status": "400", "message": err.Error()}); errJ != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}
