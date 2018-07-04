package models

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//Reprocessing handle data from discovery service
type Reprocessing struct {
	SystemID     string          `json:"systemId,omitempty"`
	PendingEvent *domain.Event   `json:"pendingEvent,omitempty"`
	Origin       *domain.Event   `json:"origin,omitempty"`
	ID           string          `json:"id,omitempty"`
	Events       []*domain.Event `json:"events,omitempty"`
	ApprovedBy   string          `json:"approvedBy,omitempty"`
}
