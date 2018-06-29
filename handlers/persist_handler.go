package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/carrot"
)

//PersistHandler handle message from persist events
func PersistHandler(context *carrot.MessageContext) error {
	eventParsed, err := getEventFromMessage(context)
	if err != nil {
		return err
	}
	bb, _ := json.Marshal(eventParsed)
	fmt.Println(string(bb))
	return context.Ack()
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
