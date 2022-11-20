package websocket

import "github.com/gorilla/websocket"

type Client struct {
	// The websocket connection.
	conn *websocket.Conn
}
