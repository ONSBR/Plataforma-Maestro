package api

import (
	"encoding/json"

	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/labstack/echo"
)

func queryReprocessing(c echo.Context) error {
	filter := map[string]string{"systemId": c.Param("systemId")}
	if c.QueryParam("status") != "" {
		filter["status"] = c.QueryParam("status")
	}
	j, err := sdk.GetDocument("reprocessing", filter)
	if err != nil {
		return err
	}
	data := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(j), &data)
	if err != nil {
		return err
	}
	for _, item := range data {
		delete(item, "_id")
	}
	return c.JSON(200, data)
}
