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