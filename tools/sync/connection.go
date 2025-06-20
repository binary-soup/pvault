package sync

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/binary-soup/go-command/util"
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

func (c Connection) SendMessage(name string, message []byte) error {
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(message)))

	_, err := c.conn.Write(length)
	if err != nil {
		return util.ChainError(err, "error writing message header to connection")
	}

	_, err = c.conn.Write(message)
	if err != nil {
		return util.ChainErrorF(err, "error writing %s to connection", name)
	}

	return nil
}

func (c Connection) ReadMessage(name string) ([]byte, error) {
	length := make([]byte, 4)

	_, err := io.ReadFull(c.conn, length)
	if err != nil {
		return nil, util.ChainError(err, "error reading message header from connection")
	}

	message := make([]byte, binary.BigEndian.Uint32(length))

	_, err = io.ReadFull(c.conn, message)
	if err != nil {
		return nil, util.ChainErrorF(err, "error reading %s from connection", name)
	}

	return message, nil
}
