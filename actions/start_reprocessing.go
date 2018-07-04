package actions

import (
	"encoding/json"
	"fmt"
	"sync"

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
		if err := mountReprocessingInfra(reprocessing); err != nil {
			log.Error(fmt.Sprintf("cannot mount reprocessing infra to system %s: ", reprocessing.SystemID), err)
			return
		}
	}
	data, _ := json.Marshal(reprocessing)
	publisher := broker.GetPublisher()
	publisher.Publish("reprocessing", reprocessing.SystemID, carrot.Message{
		ContentType: "application/json",
		Data:        data,
		Encoding:    "utf-8",
	})

	go StartReprocessing(reprocessing)
}

func mountReprocessingInfra(reprocessing models.Reprocessing) error {
	defer declareMut.Unlock()
	declareMut.Lock()
	queue := fmt.Sprintf(reprocessingQueue, reprocessing.SystemID)
	routingKey := fmt.Sprintf("#.%s.#", reprocessing.SystemID)
	if err := broker.DeclareQueue("reprocessing", queue, routingKey); err != nil {
		return err
	}
	queue = fmt.Sprintf(reprocessingEventsQueue, reprocessing.SystemID)
	routingKey = fmt.Sprintf("#.%s.#", reprocessing.SystemID)
	if err := broker.DeclareQueue("reprocessing-events", queue, routingKey); err != nil {
		return err
	}

	queue = fmt.Sprintf(reprocessingEventsControlQueue, reprocessing.SystemID)
	routingKey = fmt.Sprintf("#.%s.#", reprocessing.SystemID)
	if err := broker.DeclareQueue("reprocessing-events-control", queue, routingKey); err != nil {
		return err
	}

	queue = fmt.Sprintf(reprocessingErrorQueue, reprocessing.SystemID)
	routingKey = fmt.Sprintf("#.error_%s.#", reprocessing.SystemID)
	if err := broker.DeclareQueue("reprocessing", queue, routingKey); err != nil {
		return err
	}
	return nil
}

//StartReprocessing picks first reprocessing in queue
func StartReprocessing(reprocessing models.Reprocessing) {

	//Start reprocessing process
	context, proceed := pickReprocessing(reprocessing.ID)
	if !proceed {
		return
	}
	defer context.Nack(true)

	if err := SetStatusReprocessing(reprocessing, models.NewReprocessingStatus("running")); err != nil {
		log.Error("cannot update reprocessing on process memory: ", err)
		return
	}

	err := appdomain.PersistEntitiesByInstance(reprocessing.SystemID, reprocessing.PendingEvent.InstanceID)
	if err != nil {
		log.Error("cannot persist pending event on domain: ", err)
		if err := SetStatusReprocessing(reprocessing, models.NewReprocessingStatus("aborted:persist-domain-failure")); err != nil {
			log.Error(fmt.Sprintf("cannot set status aborted:persist-domain-failure on reprocessing %s: ", reprocessing.ID), err)
		}
		if err := context.RedirectTo("reprocessing", fmt.Sprintf("error_%s", reprocessing.ID)); err != nil {
			log.Error(fmt.Sprintf("cannot redirect reprocessing %s to error queue: ", reprocessing.ID), err)
		}
		return
	}

	if err := splitReprocessing(reprocessing); err != nil {
		log.Error(fmt.Sprintf("cannot split event for reprocessing %s: ", reprocessing.ID), err)
		if err := SetStatusReprocessing(reprocessing, models.NewReprocessingStatus("aborted:split-events-failure")); err != nil {
			log.Error(fmt.Sprintf("cannot set status aborted:persist-domain-failure on reprocessing %s: ", reprocessing.ID), err)
		}
		if err := context.RedirectTo("reprocessing", fmt.Sprintf("error_%s", reprocessing.ID)); err != nil {
			log.Error(fmt.Sprintf("cannot redirect reprocessing %s to error queue: ", reprocessing.ID), err)
		}
	}

}

func splitReprocessing(reprocessing models.Reprocessing) error {
	//TODO publish all initial events to events queue
	return nil
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
		context.Nack(true)
		return nil, false
	}
	return nil, true
}
