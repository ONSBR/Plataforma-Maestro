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

//DispatchReprocessing dispatches an approved reprocessing to system queue
func DispatchReprocessing(reprocessing models.Reprocessing, lock bool) {
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
	reprocessing.Running(lock)
	msg, _ := broker.GetMessageFrom(reprocessing)

	publisher.Publish("reprocessing", reprocessing.SystemID, msg)

}

//StartReprocessing picks first reprocessing in queue
func StartReprocessing(systemID string) {

	log.Debug("starting reprocessing")
	//Start reprocessing process
	context, proceed, reprocessing, err := pickReprocessing(systemID)
	if !proceed {
		log.Debug("reprocessing queue is empty or cannot get reprocessing")
		return
	}
	defer context.Nack(true)
	if err != nil {
		log.Error("Picking reprocessing for system: ", systemID, " has error: ", err)
		return
	}
	errFnc := func(err error) {
		if err := models.SetStatusReprocessing(reprocessing.ID, reprocessing.Status, ""); err != nil {
			log.Error(fmt.Sprintf("cannot set status aborted:persist-domain-failure on reprocessing %s: ", reprocessing.ID), err)
		}
		if err := context.RedirectTo("reprocessing", fmt.Sprintf("error_%s", reprocessing.ID)); err != nil {
			log.Error(fmt.Sprintf("cannot redirect reprocessing %s to error queue: ", reprocessing.ID), err)
		}
	}

	if err := models.SaveReprocessing(reprocessing); err != nil {
		log.Error("cannot update reprocessing on process memory: ", err)
		return
	}
	log.Debug("splitting reprocessing events")
	if err := SplitReprocessingIntoEvents(reprocessing); err != nil {
		reprocessing.AbortedSplitEventsFailure()
		log.Error(fmt.Sprintf("cannot split event for reprocessing %s: ", reprocessing.ID), err)
		errFnc(err)
		return
	}
	log.Debug("save pending commit to domain instance: ", reprocessing.PendingEvent.InstanceID)
	if err := appdomain.PersistEntitiesByInstance(reprocessing.SystemID, reprocessing.PendingEvent.InstanceID); err != nil {
		log.Error("cannot persist pending event on domain: ", err)
		reprocessing.AbortedPersistDomainFailure()
		errFnc(err)
		return
	}
	go ReprocessEvent(systemID)

}

func pickReprocessing(id string) (*carrot.MessageContext, bool, *models.Reprocessing, error) {
	picker := broker.GetPicker()
	queue := fmt.Sprintf(models.ReprocessingQueue, id)
	context, isReprocessing, err := picker.Pick(queue)
	if err != nil {
		log.Error(err)
		return nil, false, new(models.Reprocessing), err
	}
	if isReprocessing {
		reprocessing := &models.Reprocessing{}
		err := json.Unmarshal(context.Message.Data, reprocessing)
		if err != nil {
			log.Error(fmt.Sprintf("cannot unmarshall reprocessing"), err)
			return context, false, nil, err
		}
		return context, true, reprocessing, nil
	}
	return nil, false, new(models.Reprocessing), nil
}
