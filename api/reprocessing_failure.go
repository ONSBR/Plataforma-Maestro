package api

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func reprocessingFailure(c echo.Context) error {
	log.Info("reprocessing failure by event")
	event := new(domain.Event)
	if err := c.Bind(event); err != nil {
		return err
	}
	log.Info("systemId: ", event.SystemID)
	if err := actions.SetReprocessingFailure(event); err != nil {
		log.Error(err)
		return c.JSON(200, H{"message": err.Error()})
	}
	return c.JSON(200, H{"message": "ok"})
}
