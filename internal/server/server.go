package server

import (
	"context"
	"log/slog"
	"net"
	"os"
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
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.HandleClient(ctx, conn)
		}()
	}
	wg.Wait()
}
