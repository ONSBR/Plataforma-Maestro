package actions

import (
	"sync"

	"github.com/ONSBR/Plataforma-Maestro/etc"

	"github.com/PMoneda/carrot"

	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"
)

var hub map[string]chan *carrot.MessageContext
var initHub sync.Once
var mut sync.RWMutex

//ProceedToCommit process commiting to domain events by solution
func ProceedToCommit(context *carrot.MessageContext) error {
	event, _ := etc.GetEventFromMessage(context)
	if err := appdomain.PersistEntitiesByInstance(event.SystemID, event.InstanceID); err != nil {
		return err
	}
	return context.Ack()
}
