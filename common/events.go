package common

type EventName string

const (
	EventDrivers        EventName = "drivers"
	EventSession        EventName = "session"
	EventInputTelemetry EventName = "inputTelemetry"
	EventCarTelemetry   EventName = "carTelemetry"
)

type Event struct {
	Event EventName   `json:"event"`
	Data  interface{} `json:"data"`
}
