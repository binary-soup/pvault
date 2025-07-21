package sync

import (
	"fmt"
	"net"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/style"
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
		return alert.ChainError(err, "error starting host tcp server")
	}

	fmt.Printf("Listening on port %s\n", style.BoldInfo.Format(h.Port))
	return nil
}

func (h Host) Accept() (*Connection, error) {
	conn, err := h.ln.Accept()
	if err != nil {
		return nil, alert.ChainError(err, "error accepting client connection")
	}

	return NewConnection(conn), nil
}

func (h Host) Close() {
	h.ln.Close()
}
