package base

type Hub interface {
	Run()
	Register() chan Client
	Unregister() chan Client
	Comm() chan *Message
	Clients() map[Client]bool
}
