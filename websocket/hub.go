package websocket

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	clients    map[*Client]bool // Registered clients
	broadcast  chan Message     // Inbound data from the clients
	register   chan *Client     // Register requests from the clients
	unregister chan *Client     // Unregister requests from clients
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Broadcast(message Message) {
	h.broadcast <- message
}

func (h *Hub) BroadcastMessage(subscription string, data []byte) {
	h.broadcast <- NewMessage(subscription, data)
}
