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
func SetStatusReprocessing(reprocessing models.Reprocessing, status models.ReprocessingStatus) error {
	defer statusMut.Unlock()
	statusMut.Lock()

	sjson, err := sdk.GetDocument("reprocessing", map[string]string{"id": reprocessing.ID})
	if err != nil {
		return err
	}
	rep := make([]*models.Reprocessing, 1)
	json.Unmarshal([]byte(sjson), &rep)
	if len(rep) == 0 {
		return fmt.Errorf(fmt.Sprintf("no reprocessing found with id %s", reprocessing.ID))
	}
	rep[0].HistoryStatus = append(rep[0].HistoryStatus, status)
	rep[0].Status = status.Status
	return sdk.ReplaceDocument("reprocessing", map[string]string{"id": reprocessing.ID}, rep[0])
}

//GetStatusOfReprocessing return status of reprocessing on process memory
func GetStatusOfReprocessing(reprocessingID string) (string, error) {
	defer statusMut.Unlock()
	statusMut.Lock()

	sjson, err := sdk.GetDocument("reprocessing", map[string]string{"id": reprocessingID})
	if err != nil {
		return "", err
	}
	rep := make([]*models.Reprocessing, 1)
	json.Unmarshal([]byte(sjson), &rep)
	if len(rep) == 0 {
		return "", fmt.Errorf(fmt.Sprintf("no reprocessing found with id %s", reprocessingID))
	}
	return rep[0].Status, nil
}
