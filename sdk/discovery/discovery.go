package discovery

import (
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"
)

//GetReprocessingInstances from discovery service
func GetReprocessingInstances(entites appdomain.EntitiesList) ([]string, error) {
	//TODO get data from Discovey
	var en []map[string]interface{}
	apicore.Query(apicore.Filter{
		Entity: "processInstance",
		Map:    "core",
		Name:   "byProcessId",
		Params: []apicore.Param{
			apicore.Param{
				Key:   "processId",
				Value: "7828c3c7-0352-42a5-9342-2673293bc93d",
			},
		},
	}, &en)
	instances := make([]string, len(en))
	for i, obj := range en {
		instances[i] = obj["id"].(string)
	}
	return instances, nil
}
