package main

import (
	"chat-room/internal/server"
	"context"
	"log/slog"
	"os"
	"os/signal"
)

// TODO: rooms and commands
const addr string = "127.0.0.1:9000"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	
	logger := slog.Default()
	srv := server.NewServer(addr, logger)
	srv.Run(ctx)
}
