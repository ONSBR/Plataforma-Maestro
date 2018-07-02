package appdomain

import (
	"encoding/json"
	"fmt"

	"github.com/PMoneda/http"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//EntitiesList maps entities that domain app will save based on process memory
type EntitiesList []map[string]interface{}

//GetEntitiesByProcessInstance returns all entities that need to saved on domain to complete a process instance
func GetEntitiesByProcessInstance(systemID, processInstance string) (EntitiesList, *exceptions.Exception) {
	list := make(EntitiesList, 0)
	domainHost, err := getDomainHost(systemID)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/instance/%s/entities", domainHost, processInstance)

	resp, errR := http.Get(url)
	if errR != nil {
		return nil, exceptions.NewIntegrationException(errR)
	}
	errJ := json.Unmarshal([]byte(resp), &list)
	if errJ != nil {
		return nil, exceptions.NewInvalidArgumentException(errJ)
	}
	return list, nil
}

func getDomainHost(systemID string) (string, *exceptions.Exception) {
	result := make([]map[string]interface{}, 1)
	filter := apicore.Filter{
		Entity: "installedApp",
		Map:    "core",
		Name:   "bySystemIdAndType",
		Params: []apicore.Param{apicore.Param{
			Key:   "systemId",
			Value: systemID,
		}, apicore.Param{
			Key:   "type",
			Value: "domain",
		},
		},
	}
	err := apicore.Query(filter, &result)
	if err != nil {
		return "", exceptions.NewComponentException(fmt.Errorf("%s", err.Error()))
	}
	if len(result) > 0 {
		obj := result[0]
		return fmt.Sprintf("http://%s:%d", obj["host"], uint(obj["port"].(float64))), nil
	}
	return "", exceptions.NewInvalidArgumentException(fmt.Errorf("no app found for %s id", systemID))
}
