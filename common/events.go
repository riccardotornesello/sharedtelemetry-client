package common

type EventName string

const (
	VOICE_CHAT_START EventName = "VOICE_CHAT_START"
	VOICE_CHAT_END   EventName = "VOICE_CHAT_END"
)

type Event struct {
	name      EventName
	driverIdx int
}
