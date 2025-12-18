package main

import (
	"chat-room/internal/server"
	"log/slog"
	"net"
	"os"
)

const addr string = "127.0.0.1:9000"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("Listen returned error", "err", err)
		os.Exit(1)
	}

	defer listener.Close()

	slog.Info("Listening", "addr", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Accept error", "err", err)
			continue
		}
		slog.Info("Client connected", "addr", conn.RemoteAddr())
		go server.HandleClient(conn)
	}
}
