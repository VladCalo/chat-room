package server

import (
	"net"
	"log/slog"
	"strings"
	"bufio"
	"fmt"
	"context"
	"io"
)

type Client struct {
	client_id int
	name string
	conn net.Conn
	addr net.Addr
	room *Room
}

func NewClient(id int, name string, conn net.Conn) *Client {
	return &Client{
		client_id: id,
		name: name,
		conn: conn,
		addr: conn.RemoteAddr(),
	}
}

func (s *Server) addClient(conn net.Conn, name string) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++
	client := NewClient(id, name, conn)
	s.clients[id] = client

	slog.Info("Client connected", "id", client.client_id, "addr", client.addr)
	return client
}

func (s *Server) removeClient(client *Client) {
	s.mu.Lock()
	delete(s.clients, client.client_id)
	s.mu.Unlock()

	if client.room != nil {
		room := client.room
		room.mu.Lock()
		delete(room.members, client.client_id)
		room.mu.Unlock()
	}

	slog.Info("Client disconnected", "id", client.client_id, "addr", client.addr)

}

func (s *Server) HandleClient(ctx context.Context, conn net.Conn) {

	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Write([]byte("Server shutting down...Goodbye!\n"))
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	conn.Write([]byte("Enter your name: "))
	name, err := reader.ReadString('\n')
	if err != nil {
		slog.Error("Error reading client name", "err", err)
		return
	}

	client := s.addClient(conn, strings.TrimSpace(name))
	defer s.removeClient(client)

	conn.Write([]byte(fmt.Sprintf("Welcome %s!\n", client.name)))
	s.showHelp(client)

	for {
		line, err := reader.ReadString('\n')
		s.handleCommand(client, line)
		if err != nil {
			if err == io.EOF {
				return
			}
			if ctx.Err() != nil {
				slog.Info("Client connection closed", "addr", conn.RemoteAddr())
				return
			}
			slog.Error("Read error", "err", err)
			return
		}
		slog.Info("Message received", "addr", conn.RemoteAddr(), "msg", strings.TrimSpace(line))

		s.broadcast(line, client)
	}
}