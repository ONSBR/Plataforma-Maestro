package actions

import (
	"strings"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/sdk/eventmanager"
	"github.com/labstack/gommon/log"
)

func EmitErrorEvent(event *domain.Event, err error) error {
	log.Info("emitting error event with error ", err.Error())
	erroEvt := new(domain.Event)
	if strings.HasSuffix(event.Name, ".request") {
		erroEvt.Name = strings.Replace(event.Name, ".request", ".error", -1)
	} else {
		erroEvt.Name = event.Name + ".error"
	}
	erroEvt.InstanceID = event.InstanceID
	erroEvt.Branch = event.Branch
	erroEvt.IdempotencyKey = event.IdempotencyKey
	erroEvt.Image = event.Image
	erroEvt.Version = event.Version
	erroEvt.Tag = event.Tag
	erroEvt.Payload = make(map[string]interface{})
	erroEvt.Payload["instance_id"] = event.InstanceID
	erroEvt.Payload["message"] = err.Error()
	erroEvt.Payload["origin"] = event
	return eventmanager.Push(erroEvt)
}
