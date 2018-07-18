package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/gommon/log"
)

//SplitReprocessingIntoEvents takes a reprocessing e publish all events to reprocessing queue
func SplitReprocessingIntoEvents(reprocessing *models.Reprocessing) error {
	return SplitReprocessingEvents(reprocessing.ID, reprocessing.Events)
}

//SplitReprocessingEvents takes a reprocessing e publish all events to reprocessing queue
func SplitReprocessingEvents(reprocessingID string, events []*domain.Event) error {
	for i := 0; i < len(events); i++ {
		event := events[i]
		originalInstance := event.InstanceID
		event.InstanceID = ""
		event.Tag = ""
		event.Scope = "reprocessing"
		event.Reprocessing = new(domain.ReprocessingInfo)
		event.Reprocessing.ID = reprocessingID
		event.Reprocessing.InstanceID = originalInstance
		event.Reprocessing.Image = event.Image
		event.Reprocessing.SystemID = event.SystemID
		if event.Reprocessing.InstanceID == "" {
			//TODO encontrar a causa para não ter que tratar o efeito
			log.Info("exclude event from reprocessing without original instance id = ", event.Name, " branch=", event.Branch, " scope=", event.Scope, " originalInstanceID=", originalInstance)
			continue
		}
		PublishReprocessingEvents(event)
	}
	return nil
}

//PublishReprocessingEvents publish reprocessing events to reprocessing queue
func PublishReprocessingEvents(event *domain.Event) error {
	publisher := broker.GetPublisher()
	msg, _ := broker.GetMessageFrom(event)
	if err := publisher.Publish("reprocessing-events", fmt.Sprintf("%s.control_%s.backup_%s", event.SystemID, event.SystemID, event.SystemID), msg); err != nil {
		log.Error(fmt.Sprintf("failure to publish event to reprocessing-events exchange "), err)
		return err
	}
	return nil
}
