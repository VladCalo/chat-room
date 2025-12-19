package server

import (
	"net"
)

type Client struct {
	client_id int
	conn net.Conn
	addr net.Addr
}

func NewClient(id int, conn net.Conn) *Client {
	return &Client{
		client_id: id,
		conn: conn,
		addr: conn.RemoteAddr(),
	}
}