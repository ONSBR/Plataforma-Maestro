package api

import (
	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func eventGatekeeper(c echo.Context) error {
	rep, err := actions.GetReprocessingBySystemIDWithStatus(c.Param("systemId"), models.Running)
	if err != nil {
		return err
	}
	log.Debug(rep)
	if len(rep) == 0 {
		return c.NoContent(200)
	}
	return c.NoContent(403)
}
