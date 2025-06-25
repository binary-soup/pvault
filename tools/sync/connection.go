package sync

import (
	"net"
	"pvault/crypt"
	"time"
)

type Connection struct {
	conn  net.Conn
	Crypt *crypt.Crypt
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

func (c Connection) IsSecure() bool {
	return c.Crypt != nil
}

func (c Connection) SetReadTimeout(d time.Duration) {
	c.conn.SetReadDeadline(time.Now().Add(d))
}
