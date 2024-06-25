package adapter

import (
	"encoding/json"
	"io"
	"time"

	"cchalop1.com/deploy/pkg"
)

type EventServer struct {
	EventType    string    `json:"eventType"`
	Title        string    `json:"title"`
	ErrorMessage string    `json:"errorMessage"`
	Time         time.Time `json:"time"`
}

type EventServerWrapper struct {
	ServerName   string        `json:"serverName"`
	ServerId     string        `json:"serverId"`
	EventsServer []EventServer `json:"eventsServer"`
	CurrentStep  int           `json:"currentStep"`
}

func (e *EventServerWrapper) MarshalToSseEvent(w io.Writer) error {
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

func (e *EventServerWrapper) NextStep() {
	e.EventsServer[e.CurrentStep].Time = time.Now()
	e.CurrentStep += 1
}

func (e *EventServerWrapper) SetStepError(errorMessage string) {
	e.EventsServer[e.CurrentStep].ErrorMessage = errorMessage
}

type AdapterEvent struct {
	EventServerWrapper chan EventServerWrapper
}

func NewAdapterEvent() *AdapterEvent {
	return &AdapterEvent{
		EventServerWrapper: make(chan EventServerWrapper),
	}
}

func (e *AdapterEvent) SendNewEvent(event EventServerWrapper) {
	e.EventServerWrapper <- event
}
