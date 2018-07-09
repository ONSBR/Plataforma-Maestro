package actions

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/etc"
	"github.com/ONSBR/Plataforma-Maestro/sdk/eventmanager"
	"github.com/labstack/gommon/log"

	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-Maestro/models"
)

var statusMut sync.Mutex

//SetStatusReprocessing set status of reprocessing on process memory
func SetStatusReprocessing(reprocessingID string, status, owner string) error {
	rep, err := GetReprocessing(reprocessingID)
	if err != nil {
		return err
	}
	rep.SetStatus(owner, status)
	return SaveReprocessing(rep)
}

//SaveReprocessing saves reprocessing on process memory
func SaveReprocessing(reprocessing *models.Reprocessing) error {
	defer statusMut.Unlock()
	statusMut.Lock()
	return sdk.ReplaceDocument("reprocessing", map[string]string{"id": reprocessing.ID}, reprocessing)
}

//GetReprocessing return reprocessing from process memory
func GetReprocessing(reprocessingID string) (*models.Reprocessing, error) {
	defer statusMut.Unlock()
	statusMut.Lock()
	sjson, err := sdk.GetDocument("reprocessing", map[string]string{"id": reprocessingID})
	if err != nil {
		return nil, err
	}
	rep := make([]*models.Reprocessing, 1)
	err = json.Unmarshal([]byte(sjson), &rep)
	if err != nil {
		return nil, err
	}
	if len(rep) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("no reprocessing found with id %s", reprocessingID))
	}
	return rep[0], nil
}

//GetReprocessingBySystemIDWithStatus return reprocessing with systemId and status from process memory
func GetReprocessingBySystemIDWithStatus(systemID, status string) ([]*models.Reprocessing, error) {
	defer statusMut.Unlock()
	statusMut.Lock()
	sjson, err := sdk.GetDocument("reprocessing", map[string]string{"systemId": systemID, "status": status})
	if err != nil {
		return nil, err
	}
	rep := make([]*models.Reprocessing, 0)
	err = json.Unmarshal([]byte(sjson), &rep)
	return rep, err
}

//GetStatusOfReprocessing return status of reprocessing from process memory
func GetStatusOfReprocessing(reprocessingID string) (string, error) {
	rep, err := GetReprocessing(reprocessingID)
	if err != nil {
		return "", err
	}
	return rep.Status, nil
}

//ReprocessEvent takes a event to reprocess by system and emit event to event manager
func ReprocessEvent(systemID string) {
	picker := broker.GetPicker()
	context, ok, err := picker.Pick(fmt.Sprintf(reprocessingEventsQueue, systemID))
	if err != nil {
		log.Error("cannot pick event to reprocess ", err)
		return
	}
	if !ok {
		log.Error("Queue is empty")
		return
	}
	if context == nil {
		log.Error("Context is nil cannot proceed")
		return
	}
	defer context.Ack()
	event, err := etc.GetEventFromMessage(context)
	if err != nil {
		log.Error(err)
		return
	}

	err = eventmanager.Push(event)
	if err != nil {
		log.Error("cannot push event to event manager: ", err)
	}
}