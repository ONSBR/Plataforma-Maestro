package actions

import (
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//EmitEventsToReprocess put all reprocessing events to reprocessing queue
func EmitEventsToReprocess(events []*domain.Event) *exceptions.Exception {
	return nil
}
