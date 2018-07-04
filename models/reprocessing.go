package models

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/etc"
)

//Reprocessing handle data from discovery service
type Reprocessing struct {
	SystemID      string               `json:"systemId,omitempty"`
	PendingEvent  *domain.Event        `json:"pendingEvent,omitempty"`
	Origin        *domain.Event        `json:"origin,omitempty"`
	ID            string               `json:"id,omitempty"`
	Events        []*domain.Event      `json:"events,omitempty"`
	Status        string               `json:"status"`
	HistoryStatus []ReprocessingStatus `json:"history"`
}

//ReprocessingStatus stores user actions over reprocessing
type ReprocessingStatus struct {
	User      string `json:"user,omitempty"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

//NewReprocessingStatus creates a new status object to reprocessing
func NewReprocessingStatus(status string) ReprocessingStatus {
	return ReprocessingStatus{
		Status:    status,
		Timestamp: etc.GetStrTimestamp(),
	}
}
