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

// EventServerWrapper is a wrapper for the events of a server

type EventServerWrapper struct {
	ServerName       string        `json:"serverName"`
	ServerId         string        `json:"serverId"`
	EventServersList []EventServer `json:"eventServersList"`
	CurrentStep      int           `json:"currentStep"`
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
	e.EventServersList[e.CurrentStep].Time = time.Now()
	e.CurrentStep += 1
}

func (e *EventServerWrapper) SetStepError(errorMessage string) {
	e.EventServersList[e.CurrentStep].ErrorMessage = errorMessage
}

// EventDeployWrapper is a wrapper for the events of a deploy
type EventDeployWrapper struct {
	DeployName       string        `json:"deployName"`
	DeployId         string        `json:"deployId"`
	EventsDeployList []EventServer `json:"eventsDeployList"`
	CurrentStep      int           `json:"currentStep"`
}

func (e *EventDeployWrapper) MarshalToSseEvent(w io.Writer) error {
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

func (e *EventDeployWrapper) NextStep() {
	e.EventsDeployList[e.CurrentStep].Time = time.Now()
	e.CurrentStep += 1
}

func (e *EventDeployWrapper) SetStepError(errorMessage string) {
	e.EventsDeployList[e.CurrentStep].ErrorMessage = errorMessage
}

// TODO: find a better way to handle 2 types of events
type AdapterEvent struct {
	EventServerWrapper chan EventServerWrapper
	EventDeployWrapper chan EventDeployWrapper
}

func NewAdapterEvent() *AdapterEvent {
	return &AdapterEvent{
		EventServerWrapper: make(chan EventServerWrapper),
		EventDeployWrapper: make(chan EventDeployWrapper),
	}
}

func (e *AdapterEvent) SendNewServerEvent(event EventServerWrapper) {
	e.EventServerWrapper <- event
}

func (e *AdapterEvent) SendNewDeployEvent(event EventDeployWrapper) {
	e.EventDeployWrapper <- event
}
