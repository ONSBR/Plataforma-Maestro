package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"

	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/ONSBR/Plataforma-Maestro/api"
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/handlers"
	"github.com/labstack/gommon/log"
)

var local bool

const persistQueue string = "event.persist.queue"

func init() {
	flag.BoolVar(&local, "local", false, "to run service with local rabbitmq and services")
}

func main() {
	logo()

	flag.Parse()
	log.SetLevel(log.DEBUG)

	if local {
		os.Setenv("RABBITMQ_HOST", "localhost")
		os.Setenv("RABBITMQ_USERNAME", "guest")
		os.Setenv("RABBITMQ_PASSWORD", "guest")
		os.Setenv("APICORE_HOST", "localhost")
		os.Setenv("PORT", "8089")
	}

	broker.Init()
	startListenQueues()
	api.InitAPI()
}

func startListenQueues() {
	ids := getInstalledSystems()
	queues := getInputQueues(ids)
	subscribeQueues(queues)
	restartReprocessing(ids)
}

func restartReprocessing(systems []string) {
	for _, id := range systems {
		go actions.StartReprocessing(id)
	}
}

func getInstalledSystems() []string {
	type system struct {
		ID string `json:"id"`
	}
	list := make([]system, 0)
	for {
		err := apicore.Query(apicore.Filter{
			Map:    "core",
			Entity: "system",
			Name:   "",
		}, &list)
		if err == nil {
			break
		} else {
			log.Error("cannot connect to apicore retry in 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}

	ids := make([]string, len(list))
	for i, v := range list {
		ids[i] = v.ID
	}
	return ids
}

func getInputQueues(systems []string) []string {
	r := make([]string, len(systems))
	for i, id := range systems {
		r[i] = fmt.Sprintf("persist.%s.queue", id)
	}
	return r
}

func subscribeQueues(queues []string) {
	for _, q := range queues {
		err := handlers.SubscribeToReceiveEventsBySystem(q)
		if err != nil {
			panic(err)
		}
	}

}

func logo() {
	fmt.Print(`
                           _
                          | |
_ __ ___   __ _  ___   ___| |_ _ __ ___
| '_ ' _ \ / _' |/ _ \/ __| __| '__/ _ \
| | | | | | (_| |  __/\__ \ |_| | | (_) |
|_| |_| |_|\__,_|\___||___/\__|_|  \___/

`)

}
