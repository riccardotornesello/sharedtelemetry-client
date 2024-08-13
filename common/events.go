package common

type EventName string

const (
	EventDrivers        EventName = "drivers"
	EventSession        EventName = "session"
	EventInputTelemetry EventName = "inputTelemetry"
	EventCarTelemetry   EventName = "carTelemetry"
	EventRadio          EventName = "radio"
	EventFlags          EventName = "flags"
)

type Event struct {
	Event EventName   `json:"event"`
	Data  interface{} `json:"data"`
}
