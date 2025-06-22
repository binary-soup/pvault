package sync

import (
	"encoding/binary"
	"io"
	"os"
	"pvault/crypt"

	"github.com/binary-soup/go-command/util"
)

func (c Connection) ExchangeHostname() (string, error) {
	hostname, _ := os.Hostname()
	c.SendMessage("hostname", []byte(hostname))

	bytes, err := c.ReadMessage("hostname")
	return string(bytes), err
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

func (c Connection) SendSecureMessage(name string, crt *crypt.Crypt, message []byte) error {
	return c.SendMessage(name, crt.Encrypt(message))
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

func (c Connection) ReadSecureMessage(name string, crt *crypt.Crypt) ([]byte, error) {
	message, err := c.ReadMessage(name)
	if err != nil {
		return nil, err
	}

	plaintext, err := crt.Decrypt(message)
	if err != nil {
		return nil, util.ChainErrorF(err, "error decrypting %s from connection", name)
	}

	return plaintext, nil
}

func (c Connection) SendManyMessages(name string, crt *crypt.Crypt, messages [][]byte) error {
	count := make([]byte, 4)
	binary.BigEndian.PutUint32(count, uint32(len(messages)))

	_, err := c.conn.Write(count)
	if err != nil {
		return util.ChainError(err, "error writing messages count to connection")
	}

	for _, msg := range messages {
		err := c.SendSecureMessage(name, crt, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Connection) ReadManyMessages(name string, crt *crypt.Crypt, read func(uint32, uint32, []byte) error) error {
	header := make([]byte, 4)

	_, err := io.ReadFull(c.conn, header)
	if err != nil {
		return util.ChainError(err, "error reading messages count from connection")
	}
	count := binary.BigEndian.Uint32(header)

	for i := range count {
		bytes, err := c.ReadSecureMessage(name, crt)
		if err != nil {
			return err
		}

		err = read(i, count, bytes)
		if err != nil {
			return err
		}
	}
	return nil
}
