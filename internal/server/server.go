package server

import (
	"bufio"
	"context"
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
}

func NewServer(addr string, logger *slog.Logger) *Server {
	return &Server{
		addr:    addr,
		logger:  logger,
		clients: make(map[int]*Client),
		nextID:  1,
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
		client := s.addClient(conn)

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer s.removeClient(client)
			s.HandleClient(ctx, client)
		}()
	}
	wg.Wait()
}

func (s *Server) HandleClient(ctx context.Context, client *Client) {
	conn := client.conn

	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Write([]byte("Server shutting down...Goodbye!\n"))
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
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

		s.broadcast(line)
	}
}

func (s *Server) addClient(conn net.Conn) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++
	client := NewClient(id, conn)
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

func (s *Server) broadcast(msg string) {
	s.mu.Lock()
	targets := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		targets = append(targets, client)
	}
	s.mu.Unlock()

	for _, client := range targets {
		_, err := client.conn.Write([]byte(msg))
		if err != nil {
			slog.Error("Write error", "addr", client.addr, "err", err)
		}
	}
}
