package sdk

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/http"
)

//EventCanProceed to be processed by plataform
func EventCanProceed(event *domain.Event) bool {
	_, err := http.Get("")
	if err != nil {
		return false
	}
	return false
}
