package models

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/etc"
)

const Running string = "running"
const PendingApproval string = "pending_approval"
const Skipped string = "skipped"
const AbortedSplitEventsFail string = "aborted:split-events-failure"
const AbortedPersistDomainFail string = "aborted:persist-domain-failure"

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

func (rep *Reprocessing) PendingApproval() {
	rep.SetStatus("", PendingApproval)
}

func (rep *Reprocessing) IsPendingApproval() bool {
	return rep.Status == PendingApproval
}

func (rep *Reprocessing) Skipped(owner string) {
	rep.SetStatus(owner, Skipped)
}

func (rep *Reprocessing) Running() {
	rep.SetStatus("", Running)
}

func (rep *Reprocessing) AbortedSplitEventsFailure() {
	rep.SetStatus("", AbortedSplitEventsFail)
}

func (rep *Reprocessing) AbortedPersistDomainFailure() {
	rep.SetStatus("", AbortedPersistDomainFail)
}

//SetStatus on reprocessing
func (rep *Reprocessing) SetStatus(owner, status string) {
	rep.Status = status
	st := NewReprocessingStatus(rep.Status)
	st.User = owner
	if rep.HistoryStatus == nil {
		rep.HistoryStatus = make([]ReprocessingStatus, 1)
		rep.HistoryStatus[0] = st
	} else {
		rep.HistoryStatus = append(rep.HistoryStatus, st)
	}
}

func NewReprocessing(pendingEvent *domain.Event) *Reprocessing {
	return &Reprocessing{
		PendingEvent: pendingEvent,
		SystemID:     pendingEvent.SystemID,
		ID:           etc.GetUUID(),
	}
}

//NewReprocessingStatus creates a new status object to reprocessing
func NewReprocessingStatus(status string) ReprocessingStatus {
	return ReprocessingStatus{
		Status:    status,
		Timestamp: etc.GetStrTimestamp(),
	}
}
