package actions

import (
	"github.com/ONSBR/Plataforma-Maestro/models"
)

func SetStatusReprocessing(id, status string) error {
	rep, errF := models.GetReprocessing(id)
	if errF != nil {
		return errF
	}
	rep.SetStatus("user", status)
	return models.SaveReprocessing(rep)
}
