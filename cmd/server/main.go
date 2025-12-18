package main

import (
	"chat-room/internal/server"
	"log/slog"
)

// TODO: context propagation
const addr string = "127.0.0.1:9000"

func main() {
	logger := slog.Default()
	srv := server.New(addr, logger)
	srv.Run()
}
