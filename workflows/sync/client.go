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

	crt, err := w.authenticate(conn)
	if err != nil {
		return err
	}

	conn.SendSecureMessage("hostname", crt, []byte(hostname()))

	hostname, err := conn.ReadSecureMessage("hostname", crt)
	if err != nil {
		return err
	}
	printSuccessStatus(fmt.Sprintf("host identified as %s", style.BoldInfo.Format(string(hostname))))

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

		err = conn.SendMessage("header", crt.Header)
		if err != nil {
			return nil, err
		}

		status, err := conn.ReadResponse()
		if status == sync.ERROR_NONE {
			printSuccessStatus("passkey accepted")
			return crt, nil
		}
		if status == sync.ERROR_AUTH {
			printErrorStatus(err.Error())
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
