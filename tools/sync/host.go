package sync

import (
	"fmt"
	"net"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type Host struct {
	Port string
	ln   net.Listener
}

func NewHost(port string) Host {
	return Host{
		Port: port,
	}
}

func (h *Host) Start() error {
	var err error

	h.ln, err = net.Listen("tcp", h.Port)
	if err != nil {
		return util.ChainError(err, "error starting host tcp server")
	}

	fmt.Printf("Listening on port %s\n", style.BoldInfo.Format(h.Port[1:]))
	return nil
}

func (h Host) Accept() (*Connection, error) {
	conn, err := h.ln.Accept()
	if err != nil {
		return nil, util.ChainError(err, "error accepting client connection")
	}

	return NewConnection(conn), nil
}

func (h Host) Close() {
	h.ln.Close()
}
