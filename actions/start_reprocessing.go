package actions

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/PMoneda/carrot"

	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/gommon/log"
)

var cache map[string]bool
var once sync.Once
var declareMut sync.Mutex

//DispatchReprocessing dispatches an approved reprocessing to system queue
func DispatchReprocessing(reprocessing models.Reprocessing) {
	once.Do(func() {
		cache = make(map[string]bool)
	})

	if _, ok := cache[reprocessing.SystemID]; !ok {
		defer declareMut.Unlock()
		declareMut.Lock()
		queue := fmt.Sprintf("reprocessing.%s.queue", reprocessing.SystemID)
		log.Info(fmt.Sprintf("creating queue %s", queue))
		builder := broker.GetBuilder()

		if err := builder.DeclareTopicExchange("reprocessing"); err != nil {
			log.Error("Aborting reprocessing cannot declare exchange on rabbitmq: ", err)
			return
		}
		if err := builder.DeclareQueue(queue); err != nil {
			log.Error("Aborting reprocessing cannot declare a queue on rabbitmq: ", err)
			return
		}

		if err := builder.BindQueueToExchange(queue, "reprocessing", fmt.Sprintf("#.%s.#", reprocessing.SystemID)); err != nil {
			log.Error("Aborting reprocessing cannot bind queue to a exchange: ", err)
			return
		}
	}
	data, _ := json.Marshal(reprocessing)
	publisher := broker.GetPublisher()
	publisher.Publish("reprocessing", reprocessing.SystemID, carrot.Message{
		ContentType: "application/json",
		Data:        data,
		Encoding:    "utf-8",
	})
}
