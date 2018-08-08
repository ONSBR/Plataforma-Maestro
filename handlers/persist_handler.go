package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Maestro/sdk/processmemory"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/etc"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"
	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"
)

//PersistHandler handle message from persist events
func PersistHandler(context *carrot.MessageContext) (err error) {
	log.Info("received persist event")
	event, err := models.GetEventFromCeleryMessage(context)
	if err != nil {
		log.Error(err)
		context.RedirectTo("events.publish_error", "exception")
		return context.Ack()
	}
	log.Info("Processing persist by solution")
	err = handlePersistBySolution(event)
	if err == nil {
		err = context.Ack()
	} else {
		log.Error(err)
		if err := context.RedirectTo("persist_error", event.SystemID); err != nil {
			log.Error("cannot redirect to error queue: ", err)
			err = context.Nack(true)
		} else {
			err = context.Ack()
		}
	}

	return
}

func handlePersistBySolution(eventParsed *domain.Event) error {
	var err error
	origin, err := processmemory.GetEventByInstance(eventParsed.InstanceID)
	if err != nil {
		return err
	}
	if eventParsed.IsExecution() {
		log.Info("handle execution event")
		err = handleExecutionPersistence(eventParsed, origin)
	} else if eventParsed.IsReprocessing() {
		log.Info("handle reprocessing event")
		err = handleReprocessingPersistence(eventParsed)
	} else {
		err = fmt.Errorf("invalid event's scope %s", eventParsed.Scope)
	}
	if err != nil {
		log.Error(err)
		err1 := actions.EmitErrorEvent(eventParsed, err)
		if err1 != nil {
			log.Error(err1)
		}
	}
	return err
}

func handleExecutionPersistence(persistenceEvent, origin *domain.Event) (err error) {
	instances, err := actions.GetReprocessingInstances(persistenceEvent)
	instances = actions.FilterReprocessingUnit(persistenceEvent, instances)
	//log.Debug("instances to reprocess ", instances)
	events, err := actions.GetEventsFromInstances(instances)
	if err != nil {
		return err
	}
	if err == nil && len(events) > 0 {
		events = actions.SortInstances(instances, events)
		etc.LogDuration("submiting to approve reprocessing", func() {
			err = actions.SubmitReprocessingToApprove(persistenceEvent, origin, events)
		})
	} else if err == nil {
		etc.LogDuration("commiting data", func() {
			err = actions.ProceedToCommit(persistenceEvent)
		})
	}
	return
}

func handleReprocessingPersistence(eventParsed *domain.Event) (err error) {
	instances, err := actions.GetReprocessingInstances(eventParsed)
	instances = actions.FilterReprocessingUnit(eventParsed, instances)
	if err == nil && len(instances) > 0 {
		etc.LogDuration("appending reprocessing instances to reprocessing queue", func() {
			//log.Debug("get events from instances")
			events, err := actions.GetEventsFromInstances(instances)
			if err != nil {
				return
			}
			//log.Debug("get reprocessing")
			reprocessing, err := models.GetReprocessing(eventParsed.Reprocessing.ID)
			if err != nil {
				return
			}
			if !reprocessing.IsRunning() {
				err = fmt.Errorf("cannot proceed with this event because reprocessing is not running")
				return
			}
			//log.Debug("appending new reprocessing events")
			newEvents := reprocessing.AddEvents(events)
			err = models.SaveReprocessing(reprocessing)
			if err == nil {
				//log.Debug("publishing new reprocessing events")
				actions.SplitReprocessingEvents(reprocessing, newEvents)
			}
		})
	}
	if err != nil {
		log.Error("Error occurred ", err)
	}
	if err == nil {
		//log.Debug("committing data to domain")
		etc.LogDuration("commiting data", func() {
			err = appdomain.PersistEntitiesByInstance(eventParsed.SystemID, eventParsed.InstanceID)
		})
		var empty bool
		//log.Debug("pop control queue")
		_, empty, err = broker.Pop(fmt.Sprintf(models.ReprocessingEventsControlQueue, eventParsed.SystemID))
		if empty {
			err = actions.FinishReprocessing(eventParsed.SystemID)
		} else {
			//log.Debug("keep running reprocessing")
		}
	}

	if err == nil {
		//get next to execute
		//log.Debug("get next event to reprocess")
		go actions.ReprocessEvent(eventParsed.SystemID)
	}
	return
}
