package api

import "github.com/labstack/echo"

func getPendingReprocessing(c echo.Context) error {
	c.String(200, "ok")
	return nil
}
