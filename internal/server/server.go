package server

import (
	"bufio"
	"io"
	"log/slog"
	"net"
	"strings"
	"os"
)

type Server struct {
	addr string
	logger *slog.Logger
}

func New(addr string, logger *slog.Logger) *Server {
	return &Server{
		addr: addr,
		logger: logger,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		slog.Error("Listen returned error", "err", err)
		os.Exit(1)
	}

	defer listener.Close()

	slog.Info("Listening", "addr", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Accept error", "err", err)
			continue
		}
		slog.Info("Client connected", "addr", conn.RemoteAddr())
		go s.HandleClient(conn)
	}
}

func (s *Server) HandleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				slog.Info("Client disconnected", "addr", conn.RemoteAddr())
				return
			}
			slog.Error("Read error", "err", err)
			return
		}
		slog.Info("Message received", "addr", conn.RemoteAddr(), "msg", strings.TrimSpace(line))

		data := []byte(line)
		totalWritten := 0

		for totalWritten < len(data) {
			n, werr := conn.Write(data[totalWritten:])
			if werr != nil {
				slog.Error("Write error", "addr", conn.RemoteAddr(), "err", werr)
				return
			}
			totalWritten += n
		}
	}
}
