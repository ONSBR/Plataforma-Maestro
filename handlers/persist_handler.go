package handlers

import (
	"sync"
	"time"

	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/ONSBR/Plataforma-Maestro/etc"
	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"
)

var hub map[string]chan *carrot.MessageContext
var once sync.Once
var mut sync.Mutex

//PersistHandler handle message from persist events
func PersistHandler(context *carrot.MessageContext) (err error) {
	once.Do(func() {
		hub = make(map[string]chan *carrot.MessageContext)
	})
	if eventParsed, err := etc.GetEventFromMessage(context); err == nil {
		if _, ok := hub[eventParsed.SystemID]; !ok {
			mut.Lock()
			hub[eventParsed.SystemID] = make(chan *carrot.MessageContext)
			go handlePersistBySolution(hub[eventParsed.SystemID])
			mut.Unlock()
		}
		hub[eventParsed.SystemID] <- context
	}
	if err != nil {
		log.Error(err)
		err = context.RedirectTo("reprocessing_stack", "persist_error")
	}
	return
}

func handlePersistBySolution(channel chan *carrot.MessageContext) {
	for context := range channel {
		var err error
		if eventParsed, err := etc.GetEventFromMessage(context); err == nil {
			instances, err := actions.GetReprocessingInstances(eventParsed)
			if err == nil && hasReprocessing(instances) {
				err = actions.SubmitReprocessingToApprove(context, eventParsed, instances)
			} else if err == nil {
				start := time.Now()
				err = actions.ProceedToCommit(context)
				log.Info(time.Now().Sub(start))
			}
		}
		if err != nil {
			log.Error(err)
			err = context.RedirectTo("reprocessing_stack", "persist_error")
			if err != nil {
				context.Nack(true)
			}
		}
	}
}

func hasReprocessing(instances []string) bool {
	return len(instances) > 0
}
