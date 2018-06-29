package domain

//Operation is a event binding config
type Operation struct {
	Event         string `json:"event_in"`
	EventOut      string `json:"event_out"`
	Version       string `json:"version"`
	SystemID      string `json:"systemId"`
	ProcessID     string `json:"processId"`
	Reprocessable bool   `json:"reprocessable"`
}
