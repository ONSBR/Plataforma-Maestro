package discovery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"
)

//GetReprocessingInstances from discovery service
func GetReprocessingInstances(entites appdomain.EntitiesList) ([]string, error) {
	//TODO get data from Discovey
	var en []map[string]interface{}

	url := "http://process_memory:9091/instances/byEntities?systemId=ec498841-59e5-47fd-8075-136d79155705&entities=conta%2Coperacao"

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err := json.Unmarshal(body, &en)
	if err != nil {
		return nil, err
	}
	instances := make([]string, len(en))
	for i, instance := range en {
		instances[i] = instance["process"].(string)
	}
	return instances, nil
}
