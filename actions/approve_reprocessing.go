package actions

import (
	"github.com/ONSBR/Plataforma-Maestro/models"
)

//ApproveReprocessing approve reprocessing
func ApproveReprocessing(reprocessingID, user string, lock bool) (*models.Reprocessing, error) {
	reprocessing, err := models.GetReprocessingByIDWithStatus(reprocessingID, "pending_approval")
	if err != nil {
		return nil, err
	}
	reprocessings, err := models.GetManyReprocessingWithQuery(map[string]string{"tag": reprocessing.Tag, "status": "pending_approval"})
	for _, rep := range reprocessings {
		rep.Forking = !lock
		rep.Approve(user)
		err = models.SaveReprocessing(rep)
		if err != nil {
			return nil, err
		}
		DispatchReprocessing(*rep, lock)
	}
	go StartReprocessing(reprocessing.SystemID)
	return reprocessing, nil
}
