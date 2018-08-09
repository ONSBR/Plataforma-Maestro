package actions

import (
	"fmt"
	"sync"

	"github.com/PMoneda/carrot"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/eventmanager"
	"github.com/labstack/gommon/log"
)

var oneDropAtTime sync.Mutex

//DropReprocessing remove reprocessing and events from queues
func DropReprocessing(errorContext *carrot.MessageContext, systemID string) {
	oneDropAtTime.Lock()
	defer oneDropAtTime.Unlock()
	defer (func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	})()
	errorContext.Nack(true)
	context, proceed, reprocessing, err := pickReprocessing(systemID)
	if err != nil {
		log.Error(err)
		return
	}
	if !proceed {
		log.Debug("reprocessing queue is empty or cannot get reprocessing")
		return
	}
	if context == nil {
		log.Error("context is null")
		return
	}

	if reprocessing != nil {
		evt := new(domain.Event)
		evt.Name = fmt.Sprintf("%s.reprocessing.droping", systemID)
		evt.Payload = make(map[string]interface{})
		evt.Payload["reprocessing"] = reprocessing
		if err := eventmanager.Push(evt); err != nil {
			log.Error("cannot notify reprocessing drop action ", err.Error())
		}
		queue := fmt.Sprintf(models.ReprocessingEventsQueue, systemID)
		err := broker.Purge(queue)
		if err != nil {
			log.Error(err)
			context.Nack(true)
			return
		}
		queue = fmt.Sprintf(models.ReprocessingEventsControlQueue, systemID)
		if err := broker.Purge(queue); err != nil {
			log.Error(err)
			context.Nack(true)
			return
		}
		reprocessing.Failure()
		if err := models.SaveReprocessing(reprocessing); err != nil {
			log.Error(err)
			context.Nack(true)
			return
		}
	}
	context.Ack()
}
