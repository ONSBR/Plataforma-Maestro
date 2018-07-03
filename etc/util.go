package etc

import (
	"encoding/json"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/carrot"
)

//GetEventFromMessage returns an event from celery message
func GetEventFromMessage(context *carrot.MessageContext) (*domain.Event, error) {
	celeryMessage := new(domain.CeleryMessage)
	err := json.Unmarshal(context.Message.Data, celeryMessage)
	if err != nil {
		return nil, err
	}
	eventParsed := celeryMessage.Args[0]
	return &eventParsed, nil
}
