package websocket

type Message struct {
	UserId uint   `json:"userId"`
	Data   string `json:"data"`
}
