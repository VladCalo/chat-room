package main

import (
	"chat-room/internal/server"
	"fmt"
	"net"
	"os"
)

const addr string = "127.0.0.1:9000"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("[ERROR] Listen returned error", err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("[INFO] Listening on", addr)
	fmt.Println("[INFO] Waiting for one client to connect...")

	// accepting connection from clients loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[ERROR] Accept error:", err)
			continue
		}
		fmt.Println("[INFO] Client connected from:", conn.RemoteAddr())
		go server.HandleClient(conn)
	}
}
