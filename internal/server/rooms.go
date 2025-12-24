package server

import "sync"

type Room struct {
	name string
	members map[int]*Client
	mu sync.Mutex
}

func NewRoom(name string) *Room {
	return &Room{
		name: name,
		members: make(map[int]*Client),
	}
}

func (r *Room) broadcast(sender *Client, msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, client := range r.members {
		if client.client_id != sender.client_id {
			client.conn.Write([]byte(msg))
		}
	}
}