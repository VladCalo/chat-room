package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
)

type Server struct {
	addr   string
	logger *slog.Logger

	mu      sync.Mutex
	clients map[int]*Client
	nextID  int

	rooms map[string]*Room
}

func NewServer(addr string, logger *slog.Logger) *Server {
	return &Server{
		addr:    addr,
		logger:  logger,
		clients: make(map[int]*Client),
		nextID:  1,
		rooms:   make(map[string]*Room),
	}
}

func (s *Server) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		slog.Error("Listen returned error", "err", err)
		os.Exit(1)
	}

	go func() {
		<-ctx.Done()
		slog.Info("Received context cancelation (Ctrl-C)")
		listener.Close()
	}()

	slog.Info("Listening", "addr", s.addr)

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			slog.Error("Accept error", "err", err)
			continue
		}
		//client := s.addClient(conn)

		wg.Add(1)
		go func() {
			defer wg.Done()
			//defer s.removeClient(client)
			s.HandleClient(ctx, conn)
		}()
	}
	wg.Wait()
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
	defer s.mu.Unlock()

	slog.Info("Client disconnected", "id", client.client_id, "addr", client.addr)
	delete(s.clients, client.client_id)
}

func (s *Server) broadcast(msg string, client *Client) {
	s.mu.Lock()
	targets := make([]*Client, 0, len(s.clients))
	for _, c := range s.clients {
		if c.client_id != client.client_id {
			targets = append(targets, c)
		}
	}
	s.mu.Unlock()

	for _, c := range targets {
		_, err := c.conn.Write([]byte(msg))
		if err != nil {
			slog.Error("Write error", "addr", c.addr, "err", err)
		}
	}
}
