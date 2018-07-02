package actions

import (
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/carrot"
)

func SubmitReprocessingToApprove(context *carrot.MessageContext, origin *domain.Event, instances []string) *exceptions.Exception {
	return nil
}
