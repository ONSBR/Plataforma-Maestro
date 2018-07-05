package processmemory

import (
	"encoding/json"
	"fmt"

	"github.com/PMoneda/http"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//GetEventByInstance returns event from process memory
func GetEventByInstance(instanceID string) (*domain.Event, error) {
	evts := make([]*domain.Event, 0)
	resp, err := http.Get(fmt.Sprintf("%s://%s:%s/%s/event", env.Get("PROCESS_MEMORY_SCHEME", "http"), env.Get("PROCESS_MEMORY_HOST", "localhost"), env.Get("PROCESS_MEMORY_PORT", "9091"), instanceID))
	if err != nil {
		return nil, exceptions.NewIntegrationException(err)
	}
	err = json.Unmarshal([]byte(resp), &evts)
	if err != nil {
		return nil, exceptions.NewIntegrationException(err)
	}
	if len(evts) > 0 {
		return evts[0], nil
	}
	return nil, exceptions.NewInvalidArgumentException(fmt.Errorf("event not found for instance %s", instanceID))
}
