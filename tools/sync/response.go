package sync

import (
	"io"

	"github.com/binary-soup/go-command/util"
)

const (
	ERROR_NONE     = iota
	ERROR_CLIENT   = iota
	ERROR_AUTH     = iota
	ERROR_INTERNAL = iota
)

func (c Connection) SendSuccess() error {
	return c.sendResponse(ERROR_NONE, "")
}

func (c Connection) SendClientError(message string) error {
	return c.sendResponse(ERROR_CLIENT, message)
}

func (c Connection) SendAuthError(message string) error {
	return c.sendResponse(ERROR_AUTH, message)
}

func (c Connection) SendInternalError() error {
	return c.sendResponse(ERROR_INTERNAL, "internal host error")
}

func (c Connection) sendResponse(status int, message string) error {
	_, err := c.conn.Write([]byte{byte(status)})
	if err != nil {
		return util.ChainError(err, "error writing response status to connection")
	}

	if message != "" {
		return c.SendMessage("response", []byte(message))
	}
	return nil
}

func (c Connection) ReadResponse() (int, error) {
	header := make([]byte, 1)

	_, err := io.ReadFull(c.conn, header)
	if err != nil {
		return -1, util.ChainError(err, "error reading response status from connection")
	}

	status := int(header[0])

	if status == ERROR_NONE {
		return status, nil
	}

	message, err := c.ReadMessage("response")
	if err != nil {
		return status, err
	}

	return status, util.Error(string(message))
}
