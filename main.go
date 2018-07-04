package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-Maestro/api"
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/handlers"
	"github.com/PMoneda/carrot"
)

var local bool

const persistQueue string = "event.persist.queue"

func init() {
	flag.BoolVar(&local, "local", false, "to run service with local rabbitmq and services")
}

func main() {
	logo()

	flag.Parse()

	if local {
		os.Setenv("RABBITMQ_HOST", "localhost")
		os.Setenv("RABBITMQ_USERNAME", "guest")
		os.Setenv("RABBITMQ_PASSWORD", "guest")
		os.Setenv("PORT", "8089")
	}
	broker.Init()
	subscriber := broker.GetSubscriber()
	subscriber.Subscribe(carrot.SubscribeWorker{
		Queue:   persistQueue,
		Scale:   1,
		Handler: handlers.PersistHandler,
	})
	api.InitAPI()
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
