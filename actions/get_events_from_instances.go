package actions

import (
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/sdk/processmemory"
)

//GetEventsFromInstances returns all events from process memory by process instance
func GetEventsFromInstances(instances []string) ([]*domain.Event, *exceptions.Exception) {
	list := make([]*domain.Event, len(instances))
	i := 0
	for _, id := range instances {
		evt, err := processmemory.GetEventByInstance(id)
		if err != nil {
			return nil, err
		}
		list[i] = evt
		i++
	}
	return list, nil
}
