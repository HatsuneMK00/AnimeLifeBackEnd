package base

import "github.com/gorilla/websocket"

type Client interface {
	Id() uint
	Hub() Hub
	Send() chan *Message
	Conn() *websocket.Conn
}
