package actions

import (
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//SuspendCommit put current commit event to wait queue when reprocessing is executing
func SuspendCommit(event *domain.Event) *exceptions.Exception {
	//TODO
	return nil
}
