package websocket

import (
	"AnimeLifeBackEnd/global"
	"AnimeLifeBackEnd/websocket/base"
)

type hub struct {
	// Registered connections.
	clients    map[base.Client]bool
	register   chan base.Client
	unregister chan base.Client
	comm       chan *base.Message
}

func (h *hub) Register() chan base.Client {
	return h.register
}

func (h *hub) Unregister() chan base.Client {
	return h.unregister
}

func (h *hub) Comm() chan *base.Message {
	return h.comm
}

func (h *hub) Clients() map[base.Client]bool {
	return h.clients
}

func NewHub() base.Hub {
	return &hub{
		clients:    make(map[base.Client]bool),
		register:   make(chan base.Client),
		unregister: make(chan base.Client),
		comm:       make(chan *base.Message),
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			global.Logger.Infof("Client %v registered", client.Id())
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				global.Logger.Infof("Client %v unregistered", client.Id())
				delete(h.clients, client)
				close(client.Send())
			}
		case message := <-h.comm:
			global.Logger.Infof("Send message: %v", message.Data)
			// TODO send message to specific client
			for client := range h.clients {
				select {
				case client.Send() <- message:
				default:
					close(client.Send())
					delete(h.clients, client)
				}
			}
		}
	}
}
