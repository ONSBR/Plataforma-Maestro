package handlers

import (
	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/ONSBR/Plataforma-Maestro/etc"
	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"
)

//PersistHandler handle message from persist events
func PersistHandler(context *carrot.MessageContext) (err error) {

	if eventParsed, err := etc.GetEventFromMessage(context); err == nil {
		instances, err := actions.GetReprocessingInstances(eventParsed)
		if err == nil && hasReprocessing(instances) {
			err = actions.SubmitReprocessingToApprove(context, eventParsed, instances)
		} else if err == nil {
			err = actions.ProceedToCommit(context)
		}
	}
	if err != nil {
		log.Error(err)
		err = context.RedirectTo("reprocessing_stack", "persist_error")
	}
	return
}

func hasReprocessing(instances []string) bool {
	return len(instances) > 0
}
