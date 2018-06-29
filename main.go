package main

import (
	"github.com/ONSBR/Plataforma-Maestro/handlers"
	"github.com/PMoneda/carrot"
)

const persistQueue string = "event.persist.queue"

func main() {
	done := make(chan bool)
	config := carrot.ConnectionConfig{
		Host:     "localhost",
		Username: "guest",
		Password: "guest",
		VHost:    "plataforma_v1.0",
	}
	conn, _ := carrot.NewBrokerClient(&config)

	builder := carrot.NewBuilder(conn)
	builder.DeclareTopicExchange("reprocessing_stack")

	subConn, _ := carrot.NewBrokerClient(&config)

	subscriber := carrot.NewSubscriber(subConn)

	subscriber.Subscribe(carrot.SubscribeWorker{
		Queue:   persistQueue,
		Scale:   15,
		Handler: handlers.PersistHandler,
	})
	<-done
}
