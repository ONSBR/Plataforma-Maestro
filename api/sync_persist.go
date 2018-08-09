package api

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func syncPersist(c echo.Context) error {
	log.Info("sync persist started")
	event := new(domain.Event)
	if err := c.Bind(event); err != nil {
		return err
	}
	reprocessing, err := handlers.HandlePersistBySolution(event)
	if err != nil {
		return err
	}
	if reprocessing != nil {
		return c.JSON(200, H{"status": "reprocessing_request", "reprocessing": reprocessing})
	}
	return c.JSON(200, H{"status": "commited"})
}
