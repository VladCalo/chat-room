package main

import (
	"chat-room/internal/server"
	"context"
	"log/slog"
	"os"
	"os/signal"
)

// TODO: context propagation
const addr string = "127.0.0.1:9000"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	
	logger := slog.Default()
	srv := server.New(addr, logger)
	srv.Run(ctx)
}
