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

//SplitReprocessingIntoEvents takes a reprocessing e publish all events to reprocessing queue
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
		PublishReprocessingEvents(event)
	}
	return nil
}

//PublishReprocessingEvents publish reprocessing events to reprocessing queue
func PublishReprocessingEvents(event *domain.Event) error {
	log.Debug("publishing events to reprocessing events")
	publisher := broker.GetPublisher()
	log.Info("Reprocessing.InstanceID = ", event.Reprocessing.InstanceID)
	msg, _ := broker.GetMessageFrom(event)
	log.Info("Message published = ", string(msg.Data))
	if err := publisher.Publish("reprocessing-events", fmt.Sprintf("%s.control_%s", event.SystemID, event.SystemID), msg); err != nil {
		log.Error(fmt.Sprintf("failure to publish event to reprocessing-events exchange "), err)
		return err
	}
	return nil
}
