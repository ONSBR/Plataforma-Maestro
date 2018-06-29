package api

import "github.com/labstack/echo"

//reprocessingSkip skips current top reprocessing and suppress commit
func reprocessingSkip(c echo.Context) error {
	return nil
}
