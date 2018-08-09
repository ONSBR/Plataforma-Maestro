package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/eventmanager"
	"github.com/labstack/gommon/log"
)

func FinishReprocessing(systemID string) error {
	log.Debug("finalizing reprocessing")
	picker := broker.GetPicker()
	queue := fmt.Sprintf(models.ReprocessingQueue, systemID)
	context, ok, err := picker.Pick(queue)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("reprocessing queue %s is empty", queue)
	}
	if context == nil {
		return fmt.Errorf("context is nil")
	}
	defer context.Ack()
	rep, err := models.GetReprocessingFromContext(context)
	if err != nil {
		errRedi := context.RedirectTo("reprocessing", fmt.Sprintf("error_%s", systemID))
		if errRedi != nil {
			return fmt.Errorf("%s -> %s", err, errRedi)
		}
		return err
	}
	reprocessing, err := models.GetReprocessing(rep.ID)
	if err != nil {
		return err
	}
	reprocessing.Finish()
	if err := models.SaveReprocessing(reprocessing); err != nil {
		return err
	}
	evt := new(domain.Event)
	evt.Name = fmt.Sprintf("%s.reprocessing.finished", systemID)
	evt.Payload = make(map[string]interface{})
	evt.Payload["reprocessing"] = reprocessing
	go StartReprocessing(systemID)
	return eventmanager.Push(evt)
}
