package handlers

import (
	"encoding/json"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/PMoneda/carrot"
)

//PersistHandler handle message from persist events
func PersistHandler(context *carrot.MessageContext) error {
	eventParsed, err := getEventFromMessage(context)
	if err != nil {
		return err
	}
	instances, err := actions.GetReprocessingInstances(eventParsed)
	if err != nil {
		return context.Nack(true)
	}
	if hasReprocessing(instances) {
		if ex := actions.SubmitReprocessingToApprove(context, eventParsed, instances); ex != nil {
			return context.Nack(true)
		}
		return context.Ack()
	}
	if ex := actions.ProceedToCommit(eventParsed); ex != nil {
		return context.Nack(true)
	}
	return context.Ack()
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
