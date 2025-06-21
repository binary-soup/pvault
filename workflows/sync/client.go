package syncworkflow

import (
	"fmt"
	"pvault/crypt"
	"pvault/tools"
	"pvault/tools/sync"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type ClientWorkflow struct{}

func (w ClientWorkflow) Run(addr string) error {
	client := sync.NewClient()

	conn, err := client.Connect(addr)
	if err != nil {
		return nil
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", style.BoldInfo.Format(conn.RemoteAddress()))
	return w.accept(conn)
}

func (w ClientWorkflow) accept(conn *sync.Connection) error {
	hostname, err := conn.ExchangeHostname()
	if err != nil {
		return err
	}
	successLog.LogF("host identified as %s", style.BoldInfo.Format(hostname))

	crt, err := w.authenticate(conn)
	if err != nil {
		return err
	}

	conn.SendSecureMessage("test", crt, []byte("this is a test"))
	_, err = conn.ReadResponse()
	if err != nil {
		return w.hostError(err)
	}

	return nil
}

func (w ClientWorkflow) authenticate(conn *sync.Connection) (*crypt.Crypt, error) {
	for {
		passkey, err := tools.ReadPasskey("Enter Host")
		if err != nil {
			return nil, err
		}

		crt, err := crypt.NewCrypt(passkey)
		if err != nil {
			return nil, util.ChainError(err, "error creating crypt object")
		}

		conn.SendMessage("header", crt.Header)

		status, err := conn.ReadResponse()
		if status == sync.ERROR_NONE {
			successLog.Log("passkey accepted")
			return crt, nil
		}
		if status == sync.ERROR_AUTH {
			errorLog.Log(err)
			continue
		}
		if err != nil {
			return nil, w.hostError(err)
		}
	}
}

func (w ClientWorkflow) hostError(err error) error {
	return util.ChainError(err, "error from host")
}
