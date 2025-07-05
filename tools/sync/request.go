package sync

import (
	"io"

	"github.com/binary-soup/go-command/alert"
)

const (
	NO_REQUEST = 0
)

func (c Connection) SendNoRequest() error {
	return c.sendRequest("no request", NO_REQUEST, nil)
}

func (c Connection) SendRequest(name string, kind int, message []byte) error {
	return c.sendRequest(name, kind, message)
}

func (c Connection) sendRequest(name string, kind int, message []byte) error {
	_, err := c.conn.Write([]byte{byte(kind)})
	if err != nil {
		return alert.ChainError(err, "error writing request kind to connection")
	}

	if message != nil {
		return c.SendMessage(name, message)
	}
	return nil
}

func (c Connection) ReceiveRequest() (int, []byte, error) {
	header := make([]byte, 1)

	_, err := io.ReadFull(c.conn, header)
	if err != nil {
		return -1, nil, alert.ChainError(err, "error reading request kind from connection")
	}

	kind := int(header[0])

	if kind == NO_REQUEST {
		return kind, nil, nil
	}

	message, err := c.ReadMessage("response")
	return kind, message, err
}

func (c Connection) ReceiveManyRequests(handler func(int, []byte)) error {
	for {
		kind, message, err := c.ReceiveRequest()
		if err != nil {
			return err
		}

		if kind == NO_REQUEST {
			return nil
		}

		handler(kind, message)
	}
}
