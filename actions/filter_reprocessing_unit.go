package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/models"
)

//FilterReprocessingUnit filter all units by branch.
//If persistEvent is not in branch master so the filter will remove all master reprocessing units
//because a branch event cannot change master
func FilterReprocessingUnit(event *domain.Event, instances []models.ReprocessingUnit) []models.ReprocessingUnit {
	if event.Branch == "master" {
		return instances
	}
	filter := make([]models.ReprocessingUnit, 0)
	for _, ins := range instances {
		if ins.Branch == "master" {
			continue
		}
		filter = append(filter, ins)
	}
	return filter
}
