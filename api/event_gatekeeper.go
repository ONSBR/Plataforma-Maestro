package api

import (
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/echo"
)

func eventGatekeeper(c echo.Context) error {
	_, err := models.GetReprocessingBySystemIDWithStatus(c.Param("systemId"), models.Running)
	if err != nil {
		if err.Error() == "no reprocessing found" {
			return c.NoContent(200)
		}
		return err
	}
	return c.NoContent(403)
}
