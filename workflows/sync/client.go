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

func (w ClientWorkflow) Run(addr, port string) error {
	client := sync.NewClient()

	conn, err := client.Connect(addr + port)
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

	successLog.Log("receiving vault list")
	err = conn.ReadManyMessages("vault list", crt, func(bytes []byte) error {
		item, err := ParseVaultItemFromBytes(bytes)
		if err != nil {
			return util.ChainError(err, "error parsing vault item from vault list")
		}

		fmt.Printf("%s \"%s\"\n", item.ID.String(), item.Name)
		return nil
	})
	if err != nil {
		return err
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
