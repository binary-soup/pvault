package sw

import (
	"fmt"
	"passwords/crypt"
	"passwords/tools"
	"passwords/tools/sync"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

func (w SyncWorkflow) RunHost() error {
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

		retry, err := w.acceptHost(conn, passkey)
		if retry {
			fmt.Printf("  %s %s\n", style.BoldError.Format("[X]"), err.Error())
			continue
		}
		return err
	}
}

func (w SyncWorkflow) acceptHost(conn *sync.Connection, passkey string) (bool, error) {
	header, err := conn.ReadMessage("header")
	if err != nil {
		conn.SendClientError("error reading crypt header message")
		return true, err
	}

	c, invalidPasskey, err := crypt.LoadCrypt(passkey, header)
	if invalidPasskey {
		conn.SendClientError("invalid passkey")
		return true, util.Error("invalid client passkey")
	}
	if err != nil {
		conn.SendInternalError()
		return false, util.ChainError(err, "error creating crypt object")
	}

	ciphertext, err := conn.ReadMessage("ciphertext")
	if err != nil {
		conn.SendClientError("error reading ciphertext message")
		return true, err
	}

	plaintext, err := c.Decrypt(ciphertext)
	if err != nil {
		conn.SendClientError("error decrypting ciphertext")
		return true, util.ChainError(err, "error decrypting ciphertext")
	}

	fmt.Printf("%s: \"%s\"\n", style.Bolded.Format("MESSAGE"), string(plaintext))
	conn.SendSuccess()

	return false, nil
}
