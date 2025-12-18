package server

import (
	"bufio"
	"io"
	"log/slog"
	"net"
	"strings"
)

func HandleClient(conn net.Conn) {
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
