package carrot

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	client *BrokerClient
}

//Message encapsulate some data configuration
type Message struct {
	Data        []byte
	ContentType string
	Encoding    string
	Headers     map[string]interface{}
}

//Publish a message to exchange in routingkey
func (pub *Publisher) Publish(exchange, routingKey string, message Message) error {
	err := fmt.Errorf("begin")
	var ch *amqp.Channel
	for err != nil {
		if pub.client.client != nil {
			ch, err = pub.client.client.Channel()
			if err == nil {
				err = ch.Publish(
					exchange,
					routingKey,
					false,
					false,
					amqp.Publishing{
						Headers:         message.Headers,
						ContentType:     message.ContentType,
						ContentEncoding: message.Encoding,
						Body:            message.Data,
						DeliveryMode:    amqp.Persistent,
						Priority:        0,
					},
				)
				if err == nil {
					ch.Close()
				}
			}
		}
		if err != nil {
			err = pub.client.Reconnect()
			time.Sleep(1 * time.Second)
		}
	}
	return err
}

//NewPublisher creates a new broker publisher
func NewPublisher(client *BrokerClient) *Publisher {
	pub := new(Publisher)
	pub.client = client
	return pub
}
