package api

import (
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/handlers"
	"github.com/PMoneda/carrot"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func startSystemPersistHandler(c echo.Context) error {
	log.Info("starting persist handling on queue: ", c.QueryParam("queue"))
	subs := broker.GetSubscriber()
	err := subs.Subscribe(carrot.SubscribeWorker{
		Queue:   c.QueryParam("queue"),
		Scale:   10,
		Handler: handlers.PersistHandler,
	})
	if err != nil {
		return err
	}
	c.NoContent(200)
	return nil
}
