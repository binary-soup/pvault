package sync

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/binary-soup/go-command/util"
)

func (c Connection) ExchangeHostname() (string, error) {
	hostname, _ := os.Hostname()
	c.SendMessage("hostname", []byte(hostname))

	bytes, err := c.ReadMessage("hostname")
	return string(bytes), err
}

func (c Connection) SendMessage(name string, message []byte) error {
	if c.Crypt != nil {
		message = c.Crypt.Encrypt(message)
	}

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

	if c.Crypt != nil {
		message, err = c.Crypt.Decrypt(message)
		if err != nil {
			return nil, util.ChainErrorF(err, "error decrypting %s from connection", name)
		}
	}

	return message, nil
}

func (c Connection) SendMessageBlock(name string, messages [][]byte) error {
	count := make([]byte, 4)
	binary.BigEndian.PutUint32(count, uint32(len(messages)))

	_, err := c.conn.Write(count)
	if err != nil {
		return util.ChainError(err, "error writing messages count to connection")
	}

	for _, msg := range messages {
		err := c.SendMessage(name, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Connection) ReadMessageBlock(name string) ([][]byte, error) {
	count := make([]byte, 4)

	_, err := io.ReadFull(c.conn, count)
	if err != nil {
		return nil, util.ChainError(err, "error reading messages count from connection")
	}

	messages := make([][]byte, binary.BigEndian.Uint32(count))

	for i := range messages {
		bytes, err := c.ReadMessage(name)
		if err != nil {
			return nil, err
		}
		messages[i] = bytes
	}
	return messages, nil
}
