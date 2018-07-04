package broker

import (
	"fmt"
	"os"
	"time"

	"github.com/PMoneda/carrot"
)

var builder *carrot.Builder

var subscriber *carrot.Subscriber

var publisher *carrot.Publisher

//GetBuilder from current broker instance
func GetBuilder() *carrot.Builder {
	return builder
}

//GetSubscriber returns a subscriber instance
func GetSubscriber() *carrot.Subscriber {
	return subscriber
}

//GetPublisher returns publisher of broker
func GetPublisher() *carrot.Publisher {
	return publisher
}

//Init broker
func Init() {
	config := carrot.ConnectionConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
		VHost:    "plataforma_v1.0",
	}
	errC := fmt.Errorf("error")
	var conn *carrot.BrokerClient
	for errC != nil {
		conn, errC = carrot.NewBrokerClient(&config)
		time.Sleep(5 * time.Second)
	}

	builder = carrot.NewBuilder(conn)
	builder.DeclareTopicExchange("reprocessing_stack")
	builder.DeclareQueue("persist.exception_q")
	builder.DeclareQueue("create.reprocessing.exception_q")
	builder.BindQueueToExchange("persist.exception_q", "reprocessing_stack", "#.persist_error.#")
	builder.BindQueueToExchange("create.reprocessing.exception_q", "reprocessing_stack", "#.create_reprocessing_error.#")

	subConn, _ := carrot.NewBrokerClient(&config)

	subscriber = carrot.NewSubscriber(subConn)

	pubConn, _ := carrot.NewBrokerClient(&config)
	publisher = carrot.NewPublisher(pubConn)
	fmt.Println("Waiting Events")
}
