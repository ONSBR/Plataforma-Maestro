package handlers

import (
	"github.com/PMoneda/carrot"
)

//PersistHandler handle message from persist events
func PersistHandler(context *carrot.MessageContext) error {
	//Pegou a mensagem da fila de persistencia
	//Verificou se existe reprocessamento
	return context.Ack()
}
