package actions

import (
	"fmt"
	"sync"
	"time"

	"github.com/ONSBR/Plataforma-Maestro/etc"

	"github.com/PMoneda/carrot"

	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"
	"github.com/labstack/gommon/log"
)

var hub map[string]chan *carrot.MessageContext
var initHub sync.Once
var mut sync.RWMutex

//ProceedToCommit process commiting to domain events by solution
func ProceedToCommit(context *carrot.MessageContext) error {
	defer mut.Unlock()
	initHub.Do(func() {
		log.Info("Creating a channel hub")
		hub = make(map[string]chan *carrot.MessageContext)
	})
	event, err := etc.GetEventFromMessage(context)
	if err != nil {
		return err
	}
	mut.Lock()
	if _, ok := hub[event.SystemID]; !ok {
		log.Info(fmt.Sprintf("Creating a channel for solution %s", event.SystemID))
		hub[event.SystemID] = make(chan *carrot.MessageContext)
		go commit(hub[event.SystemID])
	}
	hub[event.SystemID] <- context
	return nil
}

//worker to save process result to domain
func commit(channel chan *carrot.MessageContext) {
	for context := range channel {
		event, _ := etc.GetEventFromMessage(context)
		if err := appdomain.PersistEntitiesByInstance(event.SystemID, event.InstanceID); err != nil {
			log.Error(err)
			if errR := context.RedirectTo("reprocessing_stack", "persist_error"); errR != nil {
				context.Nack(true)
				time.Sleep(10 * time.Second)
				continue
			}
		}
		if err := context.Ack(); err != nil {
			log.Error(err)
		} else {
			fmt.Printf("event %s persisted on domain\n", event.Name)
		}
	}
}
