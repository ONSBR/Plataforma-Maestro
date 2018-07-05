package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-Maestro/etc"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/processmemory"
	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"
)

//SubmitReprocessingToApprove block persistence on domain until reprocessing will be approve
func SubmitReprocessingToApprove(context *carrot.MessageContext, persistEvent *domain.Event, instances []string) (err error) {

	errorTreat := func(err error) {
		log.Error(err)
		err = context.RedirectTo("reprocessing_stack", "create_reprocessing_error")
		if err == nil {
			err = context.Ack()
		}
		return
	}
	status := "pending_approval"
	reprocessing := models.Reprocessing{
		PendingEvent:  persistEvent,
		SystemID:      persistEvent.SystemID,
		ID:            etc.GetUUID(),
		Status:        status,
		HistoryStatus: []models.ReprocessingStatus{models.ReprocessingStatus{Status: status, Timestamp: etc.GetStrTimestamp()}},
	}
	origin, err := processmemory.GetEventByInstance(persistEvent.InstanceID)
	if err != nil {
		errorTreat(err)
		return
	}
	reprocessing.Origin = origin

	events, err := getEventsFromInstances(context, instances)
	if err != nil {
		errorTreat(err)
		return
	}
	reprocessing.Events = events
	err = sdk.SaveDocument("reprocessing", reprocessing)
	if err != nil {
		errorTreat(err)
		return
	}
	err = sdk.UpdateProcessInstance(persistEvent.InstanceID, "reprocessing_pending_approval")
	if err != nil {
		errorTreat(err)
		return
	}
	context.Ack()
	return
}

func getEventsFromInstances(context *carrot.MessageContext, instances []string) ([]*domain.Event, error) {
	events := make([]*domain.Event, len(instances))
	for i, instance := range instances {
		evt, err := processmemory.GetEventByInstance(instance)
		log.Info(evt.InstanceID)
		if err != nil {
			return nil, err
		}
		events[i] = evt
	}
	return events, nil
}
