package server

import (
	"net"
)

type Client struct {
	client_id int
	name string
	conn net.Conn
	addr net.Addr
}

func NewClient(id int, name string, conn net.Conn) *Client {
	return &Client{
		client_id: id,
		name: name,
		conn: conn,
		addr: conn.RemoteAddr(),
	}
}