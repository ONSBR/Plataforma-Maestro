package discovery

import (
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"
)

//GetReprocessingInstances from discovery service
func GetReprocessingInstances(entites appdomain.EntitiesList) ([]string, *exceptions.Exception) {
	//TODO get data from Discovey
	return []string{"331231-1213123-123123"}, nil
}
