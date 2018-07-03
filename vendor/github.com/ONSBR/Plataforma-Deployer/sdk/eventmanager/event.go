package eventmanager

type Event struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}
