package actions

import (
	"sort"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/models"
)

//SortInstances organize execution order to force forking instance to be executed before other instances by branch
func SortInstances(instances []models.ReprocessingUnit, events []*domain.Event) []*domain.Event {
	e := make(domain.Events, len(events))
	for i := 0; i < len(events); i++ {
		e[i] = events[i]
	}
	sort.Sort(e)
	events = e
	find := func(id, branch string) *domain.Event {
		for i := 0; i < len(events); i++ {
			if events[i].InstanceID == id && events[i].Branch == branch {
				return events[i]
			}
		}
		return nil
	}
	finalList := make(domain.Events, len(events))
	currFork := 0
	middleList := make(domain.Events, 0)
	for i := 0; i < len(instances); i++ {
		if instances[i].Forking {
			e := find(instances[i].InstanceID, instances[i].Branch)
			finalList[currFork] = e
			currFork++
		} else {
			middleList = append(middleList, find(instances[i].InstanceID, instances[i].Branch))
		}
	}
	for i := 0; i < len(middleList); i++ {
		finalList[currFork] = middleList[i]
		currFork++
	}
	return finalList
}
