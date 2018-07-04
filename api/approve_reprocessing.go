package api

import "github.com/labstack/echo"

//ReprocessTop take the first reprocessing pending and starts to reprocess
func reprocessTop(c echo.Context) error {
	c.String(200, "ok")
	return nil
}
