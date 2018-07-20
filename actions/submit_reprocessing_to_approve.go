package actions

import (
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/gommon/log"
)

//SubmitReprocessingToApprove block persistence on domain until reprocessing will be approve
func SubmitReprocessingToApprove(persistEvent, origin *domain.Event, events []*domain.Event) (err error) {
	reprocessing := models.NewReprocessing(persistEvent)
	reprocessing.PendingApproval()

	reprocessing.Origin = origin

	reprocessing.AddEvents(events)
	err = sdk.SaveDocument("reprocessing", reprocessing)
	if err != nil {
		log.Error(err)
		return
	}
	err = sdk.UpdateProcessInstance(persistEvent.InstanceID, "reprocessing_pending_approval")
	if err != nil {
		log.Error(err)
		return
	}
	if os.Getenv(fmt.Sprintf("AUTO_REPROCESSING_%s", persistEvent.SystemID)) != "" {
		go ApproveReprocessing(reprocessing.ID, "platform", true)
	}
	list := make([]map[string]interface{}, 0)
	apicore.FindByID("processInstance", persistEvent.InstanceID, &list)
	log.Info("new reprocessing pending to approve")
	if len(list) > 0 {
		isFork, ok := list[0]["isFork"]
		if ok && isFork != nil && isFork.(bool) {
			go ApproveReprocessing(reprocessing.ID, "platform:open_branch", false)
		}
	}
	return
}
