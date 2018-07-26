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
var mutexes map[string]sync.Mutex
var onceMutexex sync.Once

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
	onceMutexex.Do(func() {
		mutexes = make(map[string]sync.Mutex)
	})
	if mutex, ok := mutexes[systemID]; !ok {
		mu := sync.Mutex{}
		mutexes[systemID] = mu
		mu.Lock()
		defer mu.Unlock()
	} else {
		mutex.Lock()
		defer mutex.Unlock()
	}
	defer (func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	})()
	log.Debug("starting reprocessing")
	//Start reprocessing process
	context, proceed, reprocessing, err := pickReprocessing(systemID)
	if !proceed {
		log.Debug("reprocessing queue is empty or cannot get reprocessing")
		return
	}
	if context == nil {
		log.Error("context is null")
		return
	}
	defer context.Nack(true)
	if err != nil {
		log.Error("Picking reprocessing for system: ", systemID, " has error: ", err)
		return
	}
	errFnc := func(err error, id, status string) {
		if reprocessing == nil {
			log.Error("cannot logging reprocessing is null")
			return
		}
		if context == nil {
			log.Error("context is null")
			return
		}
		if err := models.SetStatusReprocessing(id, status, ""); err != nil {
			log.Error(fmt.Sprintf("cannot set status aborted:persist-domain-failure on reprocessing %s: ", id), err)
			return
		}
		if err := context.RedirectTo("reprocessing", fmt.Sprintf("error_%s", id)); err != nil {
			log.Error(fmt.Sprintf("cannot redirect reprocessing %s to error queue: ", id), err)
		}
	}

	if err := models.SaveReprocessing(reprocessing); err != nil {
		log.Error("cannot update reprocessing on process memory: ", err)
		return
	}
	///log.Debug("splitting reprocessing events")
	if err := SplitReprocessingEvents(reprocessing, reprocessing.Events); err != nil {
		log.Error(fmt.Sprintf("cannot split event for reprocessing %s: ", reprocessing.ID), err)
		go DropReprocessing(context, systemID)
		reprocessing.AbortedSplitEventsFailure()
		errFnc(err, reprocessing.ID, reprocessing.Status)

		return
	}
	log.Debug("save pending commit to domain instance: ", reprocessing.PendingEvent.InstanceID)
	if err := appdomain.PersistEntitiesByInstance(reprocessing.SystemID, reprocessing.PendingEvent.InstanceID); err != nil {
		log.Error("cannot persist pending event on domain: ", err)
		go DropReprocessing(context, systemID)
		reprocessing.AbortedPersistDomainFailure()
		errFnc(err, reprocessing.ID, reprocessing.Status)
		return
	}
	go ReprocessEvent(systemID)

}

func pickReprocessing(id string) (*carrot.MessageContext, bool, *models.Reprocessing, error) {
	picker := broker.GetPicker()
	queue := fmt.Sprintf(models.ReprocessingQueue, id)
	context, isReprocessing, err := picker.Pick(queue)
	if err != nil {
		log.Error("erro on picking data on queue: ", queue)
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
