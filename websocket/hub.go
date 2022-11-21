package websocket

import "AnimeLifeBackEnd/global"

type Hub struct {
	// Registered connections.
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	comm       chan *Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			global.Logger.Infof("Client %v registered", client.id)
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				global.Logger.Infof("Client %v unregistered", client.id)
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.comm:
			// TODO send message to specific client
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
