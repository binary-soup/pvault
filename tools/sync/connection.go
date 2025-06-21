package sync

import (
	"net"
)

type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c Connection) Close() {
	c.conn.Close()
}

func (c Connection) RemoteAddress() string {
	return c.conn.RemoteAddr().String()
}
