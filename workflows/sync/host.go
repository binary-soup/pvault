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

		fmt.Printf("Connected with %s\n", style.BoldInfo.Format(conn.RemoteAddress()))

		terminate, err := w.accept(conn, passkey)
		if terminate {
			return err
		}
		errorLog.Log(err)
	}
}

func (w HostWorkflow) accept(conn *sync.Connection, passkey string) (bool, error) {
	hostname, err := conn.ExchangeHostname()
	if err != nil {
		return false, err
	}
	successLog.LogF("client identified as %s", style.BoldInfo.Format(string(hostname)))

	crt, abort, err := w.authenticate(conn, passkey)
	if err != nil {
		return abort, err
	}

	msg, err := conn.ReadSecureMessage("test", crt)
	if err != nil {
		conn.SendClientError("error reading test message")
		return false, err
	}

	successLog.Log(string(msg))
	conn.SendSuccess()

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
			errorLog.Log("invalid client passkey")
			conn.SendAuthError("invalid passkey")
			continue
		}
		if err != nil {
			conn.SendInternalError()
			return nil, true, util.ChainError(err, "error creating crypt object")
		}

		conn.SendSuccess()
		successLog.Log("client authenticated")

		return c, false, nil
	}
}
