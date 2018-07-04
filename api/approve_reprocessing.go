package api

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/ONSBR/Plataforma-Maestro/etc"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/echo"
)

type approve struct {
	User string `json:"user"`
}

var lock sync.Mutex

//approveReprocessing approve reprocessing
func approveReprocessing(c echo.Context) error {
	defer lock.Unlock()
	lock.Lock()

	approver := new(approve)
	if err := c.Bind(approver); err != nil {
		return err
	}
	j, err := sdk.GetDocument("reprocessing_pending", map[string]string{"id": c.Param("id"), "status": "pending_approval"})
	data := make([]models.Reprocessing, 0)
	json.Unmarshal([]byte(j), &data)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("reprocessing not found")
	}
	status := "approved"
	data[0].Status = status
	data[0].HistoryStatus = append(data[0].HistoryStatus, models.ReprocessingStatus{User: approver.User, Status: status, Timestamp: etc.GetStrTimestamp()})
	if err := sdk.ReplaceDocument("reprocessing_pending", map[string]string{"id": c.Param("id"), "status": "pending_approval"}, data[0]); err != nil {
		return err
	}
	go actions.DispatchReprocessing(data[0])
	return c.JSON(202, data)
}
