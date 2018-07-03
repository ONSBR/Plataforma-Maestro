package actions

import (
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/labstack/gommon/log"
)

var hub map[string]chan *domain.Event
var initHub sync.Once
var mut sync.RWMutex

//ProceedToCommit process commiting to domain events by solution
func ProceedToCommit(event *domain.Event) error {
	defer mut.Unlock()
	initHub.Do(func() {
		log.Info("Creating a channel hub")
		hub = make(map[string]chan *domain.Event)
	})
	mut.Lock()
	if _, ok := hub[event.SystemID]; !ok {
		log.Info(fmt.Sprintf("Creating a channel for solution %s", event.SystemID))
		hub[event.SystemID] = make(chan *domain.Event)
		log.Info(fmt.Sprintf("Starting a worker for solution %s", event.SystemID))
		go saveToDomain(hub[event.SystemID])
	}
	hub[event.SystemID] <- event
	return nil
}

//worker to save process result to domain
func saveToDomain(channel chan *domain.Event) {
	log.Info("Persistence worker starting to listen channel")
	for event := range channel {
		fmt.Printf("Saving event %s on domain\n", event.Name)
	}
}
