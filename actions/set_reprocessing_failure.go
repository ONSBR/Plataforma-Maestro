package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

func SetReprocessingFailure(event *domain.Event) error {
	return CleanUpFailureReprocessing(event.SystemID)
}
