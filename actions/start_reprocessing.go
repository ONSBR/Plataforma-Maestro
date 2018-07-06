package actions

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"

	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"

	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
)

var cache map[string]bool
var once sync.Once
var declareMut sync.Mutex

const reprocessingQueue = "reprocessing.%s.queue"
const reprocessingEventsQueue = "reprocessing.%s.events.queue"
const reprocessingEventsControlQueue = "reprocessing.%s.events.control.queue"
const reprocessingErrorQueue = "reprocessing.%s.error.queue"

//DispatchReprocessing dispatches an approved reprocessing to system queue
func DispatchReprocessing(reprocessing models.Reprocessing) {
	once.Do(func() {
		cache = make(map[string]bool)
	})
	if _, ok := cache[reprocessing.SystemID]; !ok {
		if err := MountReprocessingInfraBySystem(reprocessing.SystemID); err != nil {
			log.Error(fmt.Sprintf("cannot mount reprocessing infra to system %s: ", reprocessing.SystemID), err)
			return
		}
	}
	publisher := broker.GetPublisher()
	msg, _ := broker.GetMessageFrom(reprocessing)
	publisher.Publish("reprocessing", reprocessing.SystemID, msg)

	go StartReprocessing(reprocessing.SystemID)
}

//StartReprocessing picks first reprocessing in queue
func StartReprocessing(systemID string) {

	log.Debug("starting reprocessing")
	//Start reprocessing process
	context, proceed := pickReprocessing(systemID)
	if !proceed {
		log.Debug("reprocessing queue is empty")
		return
	}
	defer context.Nack(true)
	reprocessing := models.Reprocessing{}
	err := json.Unmarshal(context.Message.Data, &reprocessing)
	if err != nil {
		log.Error(fmt.Sprintf("cannot unmarshall reprocessing for system %s: ", systemID), err)
		return
	}

	log.Debug("set reprocessing to running")
	if err := SetStatusReprocessing(reprocessing, models.NewReprocessingStatus("running")); err != nil {
		log.Error("cannot update reprocessing on process memory: ", err)
		return
	}
	log.Debug("save pending commit to domain")
	if err := appdomain.PersistEntitiesByInstance(reprocessing.SystemID, reprocessing.PendingEvent.InstanceID); err != nil {
		log.Error("cannot persist pending event on domain: ", err)
		if err := SetStatusReprocessing(reprocessing, models.NewReprocessingStatus("aborted:persist-domain-failure")); err != nil {
			log.Error(fmt.Sprintf("cannot set status aborted:persist-domain-failure on reprocessing %s: ", reprocessing.ID), err)
		}
		if err := context.RedirectTo("reprocessing", fmt.Sprintf("error_%s", reprocessing.ID)); err != nil {
			log.Error(fmt.Sprintf("cannot redirect reprocessing %s to error queue: ", reprocessing.ID), err)
		}
		return
	}
	log.Debug("splitting reprocessing events")

	if err := SplitReprocessingIntoEvents(reprocessing); err != nil {
		log.Error(fmt.Sprintf("cannot split event for reprocessing %s: ", reprocessing.ID), err)
		if err := SetStatusReprocessing(reprocessing, models.NewReprocessingStatus("aborted:split-events-failure")); err != nil {
			log.Error(fmt.Sprintf("cannot set status aborted:persist-domain-failure on reprocessing %s: ", reprocessing.ID), err)
		}
		if err := context.RedirectTo("reprocessing", fmt.Sprintf("error_%s", reprocessing.ID)); err != nil {
			log.Error(fmt.Sprintf("cannot redirect reprocessing %s to error queue: ", reprocessing.ID), err)
		}
	}

}

func pickReprocessing(id string) (*carrot.MessageContext, bool) {
	picker := broker.GetPicker()
	queue := fmt.Sprintf(reprocessingQueue, id)
	context, isReprocessing, err := picker.Pick(queue)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	if isReprocessing {
		//give back reprocessing message to queue and proceed
		return context, true
	}
	return nil, false
}

func getOperationFromEvent(event *domain.Event) (*domain.OperationInstance, error) {
	log.Info(event.InstanceID, " ", event.Name)
	operations := make([]*domain.OperationInstance, 0)
	err := apicore.Query(apicore.Filter{
		Entity: "operationInstance",
		Map:    "core",
		Name:   "byInstanceIdEventName",
		Params: []apicore.Param{
			apicore.Param{
				Key:   "processInstanceId",
				Value: event.InstanceID,
			},
			apicore.Param{
				Key:   "eventName",
				Value: event.Name,
			},
		},
	}, &operations)
	if err != nil {
		return nil, err
	}
	if len(operations) > 0 {
		return operations[0], nil
	}
	return nil, fmt.Errorf(fmt.Sprintf("operation instance for event %s and instance %s not found", event.Name, event.InstanceID))
}
