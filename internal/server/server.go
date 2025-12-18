package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func HandleClient(conn net.Conn) {
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
