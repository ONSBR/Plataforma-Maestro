package actions

import (
	"encoding/json"
	"fmt"
	"sync"

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
	json.Unmarshal([]byte(sjson), &rep)
	if len(rep) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("no reprocessing found with id %s", reprocessingID))
	}
	return rep[0], nil
}

//GetStatusOfReprocessing return status of reprocessing from process memory
func GetStatusOfReprocessing(reprocessingID string) (string, error) {
	rep, err := GetReprocessing(reprocessingID)
	if err != nil {
		return "", err
	}
	return rep.Status, nil
}
