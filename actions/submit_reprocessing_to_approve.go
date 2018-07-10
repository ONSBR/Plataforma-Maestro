package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/processmemory"
	"github.com/labstack/gommon/log"
)

//SubmitReprocessingToApprove block persistence on domain until reprocessing will be approve
func SubmitReprocessingToApprove(persistEvent *domain.Event, instances []string) (err error) {
	reprocessing := models.NewReprocessing(persistEvent)
	reprocessing.PendingApproval()
	origin, err := processmemory.GetEventByInstance(persistEvent.InstanceID)
	if err != nil {
		log.Error(err)
		return
	}
	reprocessing.Origin = origin

	events, err := getEventsFromInstances(instances)
	if err != nil {
		log.Error(err)
		return
	}
	reprocessing.Events = events
	err = sdk.SaveDocument("reprocessing", reprocessing)
	if err != nil {
		log.Error(err)
		return
	}
	err = sdk.UpdateProcessInstance(persistEvent.InstanceID, "reprocessing_pending_approval")
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func getEventsFromInstances(instances []string) ([]*domain.Event, error) {
	events := make([]*domain.Event, len(instances))
	for i, instance := range instances {
		evt, err := processmemory.GetEventByInstance(instance)
		if err != nil {
			return nil, err
		}
		events[i] = evt
	}
	return events, nil
}
