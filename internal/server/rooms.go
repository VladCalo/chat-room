package server

type Room struct {
	name string
	members map[int]*Client
}

func NewRoom(name string) *Room {
	return &Room{
		name: name,
		members: make(map[int]*Client),
	}
}