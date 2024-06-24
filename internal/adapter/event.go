package adapter

import (
	"encoding/json"
	"io"
	"time"

	"cchalop1.com/deploy/pkg"
)

type EventServer struct {
	EventType    string    `json:"event_type"`
	ServerId     string    `json:"server_id"`
	Message      string    `json:"message"`
	ErrorMessage *string   `json:"error_message"`
	Step         int       `json:"step"`
	TotalSteps   int       `json:"total_steps"`
	Time         time.Time `json:"time"`
}

func (e *EventServer) MarshalToSseEvent(w io.Writer) error {
	Data, err := json.Marshal(e)

	if err != nil {
		return err
	}

	event := pkg.SseEvent{
		Data: Data,
	}

	event.MarshalTo(w)

	return nil
}

type AdapterEvent struct {
	EventServer chan EventServer
}

func NewAdapterEvent() *AdapterEvent {
	return &AdapterEvent{
		EventServer: make(chan EventServer),
	}
}

func (e *AdapterEvent) CreateNewEvent(event EventServer) {

	e.EventServer <- event
}
