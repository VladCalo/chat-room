package server

import (
	"fmt"
	"strings"
)

func (s *Server) handleCommand(client *Client, input string) {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "/join":
		if len(parts) < 2 {
			s.sendToClient(client, "Usage: /join <room>\n")
			return
		}
		s.joinRoom(parts[1], client)
	case "/exit":
	case "/list":
	case "/whereAmI":
		s.whereAmI(client)
	case "/help":
		s.showHelp(client)
	default:
		s.sendToClient(client, "Unknown command. Type /help\n")
	}
}

func (s *Server) showHelp(client *Client) {
	help := `
#######################################################
#  Commands:                                          #
#    /join <room>  - Join a room                      #
#    /exit         - Leave current room               #
#    /list         - List all rooms                   #
#    /whereAmI     - Tells you which room you are in  #
#    /help         - Show this help                   #
#######################################################
`
	s.sendToClient(client, help)
}

func (s *Server) sendToClient(client *Client, msg string) {
	client.conn.Write([]byte(msg))
}

func (s *Server) joinRoom(roomName string, client *Client) {
	if client.room != nil && client.room.name == roomName {
		s.sendToClient(client, "Your are already in this room!\n")
		return
	}

	s.mu.Lock()
	room, exists := s.rooms[roomName]

	if !exists {
		room = NewRoom(roomName)
		s.rooms[roomName] = room
	}
	s.mu.Unlock()

	room.mu.Lock()
	room.members[client.client_id] = client
	room.mu.Unlock()

	client.room = room

	s.sendToClient(client, fmt.Sprintf("Joined room: %s\n", roomName))
}

//func (s *Server) exitRoom(client)

func (s *Server) whereAmI(client *Client) {
	if client.room == nil {
		s.sendToClient(client, "You are in the lobby...Please join a room to send messages!\n")
		return
	}
	s.sendToClient(client, fmt.Sprintf("You are in room: %s\n", client.room.name))
}