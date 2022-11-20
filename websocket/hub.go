package websocket

type Hub struct {
	// Registered connections.
	connections map[*Client]bool
}
