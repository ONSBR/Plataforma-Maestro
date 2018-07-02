package appdomain

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//EntitiesList maps entities that domain app will save based on process memory
type EntitiesList []DomainEntity

type DomainEntity map[string]interface{}

//GetEntitiesByProcessInstance returns all entities that need to saved on domain to complete a process instance
func GetEntitiesByProcessInstance(systemID, processInstance string) (EntitiesList, *exceptions.Exception) {
	list := make(EntitiesList, 0)
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
	err := apicore.Query(filter, &list)
	if err != nil {
		return nil, exceptions.NewComponentException(fmt.Errorf("%s", err.Error()))
	}
	return list, nil
}
