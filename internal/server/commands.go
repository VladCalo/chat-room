package server

import (
	"strings"
)

func (s *Server) handleCommand(client *Client, input string) {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "/join":
	case "/exit":
	case "/list":
	case "/help":
		s.showHelp(client)
	default:
		s.sendToClient(client, "Unknown command. Type /help\n")
	}
}

func (s *Server) showHelp(client *Client) {
	help := `
####################################################
#  Commands:                                       #
#    /join <room>  - Join a room                   #
#    /exit         - Leave current room            #
#    /list         - List all rooms                #
#    /help         - Show this help                #
####################################################
`
	s.sendToClient(client, help)
}

func (s *Server) sendToClient(client *Client, msg string) {
	client.conn.Write([]byte(msg))
}
