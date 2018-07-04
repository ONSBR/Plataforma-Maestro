package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-Maestro/api"
	"github.com/ONSBR/Plataforma-Maestro/handlers"
	"github.com/PMoneda/carrot"
)

const persistQueue string = "event.persist.queue"

var local bool

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
	config := carrot.ConnectionConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
		VHost:    "plataforma_v1.0",
	}
	conn, _ := carrot.NewBrokerClient(&config)

	builder := carrot.NewBuilder(conn)
	builder.DeclareTopicExchange("reprocessing_stack")
	builder.DeclareQueue("persist.exception_q")
	builder.DeclareQueue("create.reprocessing.exception_q")
	builder.BindQueueToExchange("persist.exception_q", "reprocessing_stack", "#.persist_error.#")
	builder.BindQueueToExchange("create.reprocessing.exception_q", "reprocessing_stack", "#.create_reprocessing_error.#")

	subConn, _ := carrot.NewBrokerClient(&config)

	subscriber := carrot.NewSubscriber(subConn)

	subscriber.Subscribe(carrot.SubscribeWorker{
		Queue:   persistQueue,
		Scale:   1,
		Handler: handlers.PersistHandler,
	})
	fmt.Println("Waiting Events")
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
