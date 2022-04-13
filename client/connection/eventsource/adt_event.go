package eventsource

import (
	"encoding/json"
)

type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Data any    `json:"data"`
}

func (this *Event) ToBytes() []byte {
	data, err := json.Marshal(this.Data)

	if err != nil {
		panic(err)
	}

	event := ""

	if this.ID != "" {
		event += "id: " + this.ID + "\n"
	}

	event += "event: " + this.Name + "\ndata: "

	return append([]byte(event), data...)
}

func (this *Event) ToString() string {
	return string(this.ToBytes())
}
