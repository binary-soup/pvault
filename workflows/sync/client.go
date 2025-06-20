package sw

import (
	"fmt"
	"passwords/crypt"
	"passwords/tools"
	"passwords/tools/sync"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

func (w SyncWorkflow) RunClient(addr string) error {
	client := sync.NewClient()

	conn, err := client.Connect(addr)
	if err != nil {
		return nil
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", style.BoldInfo.Format(conn.RemoteAddress()))

	c, err := w.authenticateWithHost(conn)
	if err != nil {
		return err
	}

	ciphertext := c.Encrypt([]byte("This is a test message"))

	err = conn.SendMessage("ciphertext", ciphertext)
	if err != nil {
		return nil
	}

	_, err = conn.ReadResponse()
	if err != nil {
		return w.newHostError(err)
	}

	w.printSuccessStatus("message received")
	return nil
}

func (w SyncWorkflow) authenticateWithHost(conn *sync.Connection) (*crypt.Crypt, error) {
	for {
		passkey, err := tools.ReadPasskey("Enter Host")
		if err != nil {
			return nil, err
		}

		c, err := crypt.NewCrypt(passkey)
		if err != nil {
			return nil, util.ChainError(err, "error creating crypt object")
		}

		err = conn.SendMessage("header", c.Header)
		if err != nil {
			return nil, err
		}

		status, err := conn.ReadResponse()
		if status == sync.ERROR_NONE {
			w.printSuccessStatus("passkey accepted")
			return c, nil
		}
		if status == sync.ERROR_AUTH {
			w.printErrorStatus(err.Error())
			continue
		}
		if err != nil {
			return nil, w.newHostError(err)
		}
	}
}

func (w SyncWorkflow) newHostError(err error) error {
	return util.ChainError(err, "error from host")
}
