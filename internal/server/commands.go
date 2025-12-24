package server

import (
	"fmt"
	"strings"
)

const lobbyMsg string = "You are in the lobby...Please join a room to send messages!\n"

func (s *Server) handleCommand(client *Client, input string) {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "/join":
		if len(parts) != 2 {
			s.sendToClient(client, "Usage: /join <room>\n")
			return
		}
		s.joinRoom(parts[1], client)
	case "/exit":
		if len(parts) != 1 {
			s.sendToClient(client, "Usage: /exit\n")
			return
		}
		s.exitRoom(client)
	case "/list":
		if len(parts) == 1 {
			s.listRooms(client)
		} else if len(parts) == 2 {
			s.listRoomMembers(client, parts[1])
		} else {
			s.sendToClient(client, "Usage: /list or /list <room>\n")
		}
	case "/whereAmI":
		if len(parts) != 1 {
			s.sendToClient(client, "Usage: /whereAmI\n")
			return
		}
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
#    /list <room>  - List people in specified room     #
#    /whereAmI     - Tells you which room you are in  #
#    /help         - Show this help                   #
#    Control-C     - Disconnect                       #
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

	//exit current room if trying to connect to another room already being in a room
	if client.room != nil {
		s.exitRoom(client)
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

func (s *Server) exitRoom(client *Client) {
	if client.room == nil {
		s.sendToClient(client, lobbyMsg)
		return
	}

	room := client.room

	room.mu.Lock()
	delete(room.members, client.client_id)
	room.mu.Unlock()

	client.room = nil
	s.sendToClient(client, fmt.Sprintf("Left room: %s\n", room.name))
}

func (s *Server) listRooms(client *Client) {
	s.mu.Lock()
	if len(s.rooms) == 0 {
		s.mu.Unlock()
		s.sendToClient(client, "Rooms: None. Joining a room will automatically create one!\n")
		return
	}

	var sb strings.Builder
	sb.WriteString("Rooms:\n")

	for roomName := range s.rooms {
		sb.WriteString("  " + roomName + "\n")
	}
	output := sb.String()
	s.mu.Unlock()

	s.sendToClient(client, output)
}

func (s *Server) listRoomMembers(client *Client, roomName string) {
	s.mu.Lock()
	room, exists := s.rooms[roomName]
	s.mu.Unlock()

	if !exists {
		s.sendToClient(client, fmt.Sprintf("Room '%s' does not exist\n", roomName))
		return
	}

	room.mu.Lock()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Room %s:\n", room.name))

	for _, c := range room.members {
		sb.WriteString(fmt.Sprintf("  - %s\n", c.name))
	}
	output := sb.String()
	room.mu.Unlock()

	s.sendToClient(client, output)
}

func (s *Server) whereAmI(client *Client) {
	if client.room == nil {
		s.sendToClient(client, lobbyMsg)
		return
	}
	s.sendToClient(client, fmt.Sprintf("You are in room: %s\n", client.room.name))
}
