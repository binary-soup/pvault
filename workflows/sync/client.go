package sw

import (
	"passwords/crypt"
	"passwords/tools"
	"passwords/tools/sync"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

func (w SyncWorkflow) RunClient(addr string) error {
	const MESSAGE = "This is a test message."

	passkey, err := tools.ReadPasskey("Enter Host")
	if err != nil {
		return err
	}

	client := sync.NewClient()

	conn, err := client.Connect(addr)
	if err != nil {
		return nil
	}
	defer conn.Close()

	c, err := crypt.NewCrypt(passkey)
	if err != nil {
		return util.ChainError(err, "error creating crypt object")
	}

	ciphertext := c.Encrypt([]byte(MESSAGE))

	err = conn.SendMessage("header", c.Header)
	if err != nil {
		return nil
	}

	err = conn.SendMessage("ciphertext", ciphertext)
	if err != nil {
		return nil
	}

	err = conn.ReadResponse()
	if err != nil {
		return util.ChainError(err, "error from host")
	}

	style.BoldSuccess.Println("Success!")
	return nil
}
