package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"net"
	"os"
)

const addr string = "127.0.0.1:9000"

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	//reading loop for current conn
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("[INFO] Client disconnected:", conn.RemoteAddr())
				return
			}
			fmt.Println("[ERROR] Read had error:", err)
			return
		}
		line = strings.TrimSpace(line)
		fmt.Printf("[INFO] From %v: %q\n", conn.RemoteAddr(), line)

		data := []byte(line)
		totalWritten := 0
		
		// write data back to client
		for totalWritten < len(data) {
			n, werr := conn.Write(data[totalWritten:])
			if werr != nil {
				fmt.Println("[ERROR] Writing error to", conn.RemoteAddr(), ":", werr)
				return
			}
			totalWritten += n
		}
	}
}

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
		go handleClient(conn)
	}
}