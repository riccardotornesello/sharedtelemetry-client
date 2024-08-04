package websocket

type Message struct {
	subscription string
	data         []byte
}

func NewMessage(subscription string, data []byte) Message {
	return Message{
		subscription: subscription,
		data:         data,
	}
}
