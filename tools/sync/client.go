package sync

import (
	"net"

	"github.com/binary-soup/go-command/alert"
)

type Client struct{}

func NewClient() Client {
	return Client{}
}

func (c Client) Connect(address string) (*Connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, alert.ChainErrorF(err, "error dialing host at \"%s\"", address)
	}

	return NewConnection(conn), nil
}
