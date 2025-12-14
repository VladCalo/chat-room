package main

import (
	"fmt"
	"io"
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

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("[EROR] Error accepting connections:", err)
		os.Exit(1)
	}

	defer conn.Close()
	
	fmt.Println("[INFO] Client connected from", conn.RemoteAddr())
	fmt.Println("[INFO] Listen loop started...")

	buf := make([]byte, 4)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client disconnected")
				return
			}
			fmt.Println("[ERROR] Read error:", err)
			os.Exit(1)
		}

		received := buf[:n]
		fmt.Printf("Received %d bytes: %q\n", n, string(received))

		totalWritten := 0
		for totalWritten < n {
			w, werr := conn.Write(received[totalWritten:])
			if werr != nil {
				fmt.Println("[ERROR] Write error:", werr)
				return
			}
			totalWritten += w
		}
	}

}