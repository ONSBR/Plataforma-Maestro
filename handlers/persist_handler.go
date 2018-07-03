package handlers

import (
	"encoding/json"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"
)

//PersistHandler handle message from persist events
func PersistHandler(context *carrot.MessageContext) error {
	eventParsed, err := getEventFromMessage(context)
	if err != nil {
		return err
	}
	errAck := context.Ack()
	instances, err := actions.GetReprocessingInstances(eventParsed)
	if err != nil {
		log.Error(err)
		return context.RedirectTo("reprocessing_stack", "persist_error")
	}
	if hasReprocessing(instances) {
		if ex := actions.SubmitReprocessingToApprove(context, eventParsed, instances); ex != nil {
			log.Error(err)
			return context.RedirectTo("reprocessing_stack", "persist_error")
		}
	}
	if ex := actions.ProceedToCommit(eventParsed); ex != nil {
		log.Error(err)
		return context.RedirectTo("reprocessing_stack", "persist_error")
	}
	return errAck
}

func hasReprocessing(instances []string) bool {
	return len(instances) > 0
}

func getEventFromMessage(context *carrot.MessageContext) (*domain.Event, error) {
	celeryMessage := new(domain.CeleryMessage)
	err := json.Unmarshal(context.Message.Data, celeryMessage)
	if err != nil {
		return nil, err
	}
	eventParsed := celeryMessage.Args[0]
	return &eventParsed, nil
}
