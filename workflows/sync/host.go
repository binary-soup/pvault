package syncworkflow

import (
	"fmt"
	"pvault/crypt"
	"pvault/tools"
	"pvault/tools/sync"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type HostWorkflow struct{}

func (w HostWorkflow) Run() error {
	passkey, err := tools.ReadAndVerifyPasskey("Choose Host")
	if err != nil {
		return err
	}

	host := sync.NewHost(":9000")

	err = host.Start()
	if err != nil {
		return err
	}
	defer host.Close()

	for {
		conn, err := host.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()

		terminate, err := w.accept(conn, passkey)
		if terminate {
			return err
		}
		printErrorStatus(err.Error())
	}
}

func (w HostWorkflow) accept(conn *sync.Connection, passkey string) (bool, error) {
	fmt.Printf("Connected with %s\n", style.BoldInfo.Format(conn.RemoteAddress()))

	crt, abort, err := w.authenticate(conn, passkey)
	if abort {
		return true, err
	}

	conn.SendSecureMessage("hostname", crt, []byte(hostname()))

	hostname, err := conn.ReadSecureMessage("hostname", crt)
	if err != nil {
		conn.SendClientError("error reading hostname message")
		return false, err
	}
	printSuccessStatus(fmt.Sprintf("client identified as %s", style.BoldInfo.Format(string(hostname))))

	return true, nil
}

func (w HostWorkflow) authenticate(conn *sync.Connection, passkey string) (*crypt.Crypt, bool, error) {
	var c *crypt.Crypt
	var invalidPasskey bool

	for {
		header, err := conn.ReadMessage("header")
		if err != nil {
			conn.SendClientError("error reading crypt header message")
			return nil, false, err
		}

		c, invalidPasskey, err = crypt.LoadCrypt(passkey, header)
		if invalidPasskey {
			printErrorStatus("invalid client passkey")
			conn.SendAuthError("invalid passkey")
			continue
		}
		if err != nil {
			conn.SendInternalError()
			return nil, true, util.ChainError(err, "error creating crypt object")
		}

		printSuccessStatus("client authenticated")
		conn.SendSuccess()
		return c, false, nil
	}
}
