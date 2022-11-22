package websocket

import (
	"AnimeLifeBackEnd/global"
	"AnimeLifeBackEnd/websocket/base"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	pingPeriod = 45 * time.Second
)

type client struct {
	// The websocket connection.
	id   uint
	hub  base.Hub
	conn *websocket.Conn
	send chan *base.Message
}

func (c *client) Id() uint {
	return c.id
}

func (c *client) Hub() base.Hub {
	return c.hub
}

func (c *client) Send() chan *base.Message {
	return c.send
}

func (c *client) Conn() *websocket.Conn {
	return c.conn
}

func (c *client) readPump() {
	defer func() {
		c.hub.Unregister() <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	for {
		// once read message failed, the connection will be closed
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Errorf("ReadMessage error: %v", err)
			}
			break
		}
		global.Logger.Infof("Send to Client %v message: %v", c.id, message)
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.conn.Close()
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				global.Logger.Warnf("The hub closed the channel for client %v", c.id)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteJSON(message)
			if err != nil {
				global.Logger.Errorf("WriteJson error: %v", err)
				return
			}

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				if err := c.conn.WriteJSON(<-c.send); err != nil {
					global.Logger.Errorf("WriteJson error: %v", err)
					return
				}
			}
		case <-ticker.C:
			global.Logger.Infof("Client %v ping", c.id)
			err := c.conn.WriteJSON(&base.Message{Type: "ping"})
			if err != nil {
				global.Logger.Errorf("WriteJson error: %v", err)
				return
			}
		}
	}
}

func ServeWs(hub base.Hub, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		global.Logger.Errorf("Upgrade error: %v", err)
		return
	}

	client := &client{hub: hub, conn: conn, send: make(chan *base.Message, 256)}
	client.hub.Register() <- client

	go client.writePump()
	go client.readPump()

	//defer conn.Close()
	//
	//for {
	//	var message Message
	//	err := conn.ReadJSON(&message)
	//	if err != nil {
	//		global.Logger.Errorf("ReadJson error: %v", err)
	//		break
	//	}
	//	global.Logger.Infof("message: %v", message.Data)
	//}
}
