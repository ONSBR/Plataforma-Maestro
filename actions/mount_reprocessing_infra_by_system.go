package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Maestro/broker"
)

//MountReprocessingInfraBySystem mounts rabbitmq infra to handle reprocessing by system
func MountReprocessingInfraBySystem(systemID string) error {
	defer declareMut.Unlock()
	declareMut.Lock()
	queue := fmt.Sprintf(reprocessingQueue, systemID)
	routingKey := fmt.Sprintf("#.%s.#", systemID)
	if err := broker.DeclareQueue("reprocessing", queue, routingKey); err != nil {
		return err
	}
	queue = fmt.Sprintf(reprocessingEventsQueue, systemID)
	routingKey = fmt.Sprintf("#.%s.#", systemID)
	if err := broker.DeclareQueue("reprocessing-events", queue, routingKey); err != nil {
		return err
	}

	queue = fmt.Sprintf(reprocessingEventsControlQueue, systemID)
	routingKey = fmt.Sprintf("#.control_%s.#", systemID)
	if err := broker.DeclareQueue("reprocessing-events", queue, routingKey); err != nil {
		return err
	}

	queue = fmt.Sprintf(reprocessingErrorQueue, systemID)
	routingKey = fmt.Sprintf("#.error_%s.#", systemID)
	if err := broker.DeclareQueue("reprocessing", queue, routingKey); err != nil {
		return err
	}
	return nil
}
