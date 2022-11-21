package base

type Message struct {
	UserId uint   `json:"userId"`
	Type   string `json:"type"`
	Data   string `json:"data"`
}
