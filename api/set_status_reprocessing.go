package api

import (
	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/labstack/echo"
)

func setStatusReprocessing(c echo.Context) error {
	id := c.QueryParam("id")
	status := c.QueryParam("status")
	return actions.SetStatusReprocessing(id, status)
}
