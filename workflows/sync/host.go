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
			w.printErrorStatus(err.Error())
			continue
		}
		return err
	}
}

func (w SyncWorkflow) acceptHost(conn *sync.Connection, passkey string) (bool, error) {
	fmt.Printf("Connected with %s\n", style.BoldInfo.Format(conn.RemoteAddress()))

	c, retry, err := w.authenticateClient(conn, passkey)
	if retry {
		return true, err
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

	w.printSuccessStatus(fmt.Sprintf("%s: \"%s\"", style.Bolded.Format("MESSAGE"), string(plaintext)))
	conn.SendSuccess()

	return false, nil
}

func (w SyncWorkflow) authenticateClient(conn *sync.Connection, passkey string) (*crypt.Crypt, bool, error) {
	var c *crypt.Crypt
	var invalidPasskey bool

	for {
		header, err := conn.ReadMessage("header")
		if err != nil {
			conn.SendClientError("error reading crypt header message")
			return nil, true, err
		}

		c, invalidPasskey, err = crypt.LoadCrypt(passkey, header)
		if invalidPasskey {
			w.printErrorStatus("invalid client passkey")
			err := conn.SendAuthError("invalid passkey")
			if err != nil {
				return nil, true, util.ChainError(err, "error sending auth error")
			}
			continue
		}
		if err != nil {
			conn.SendInternalError()
			return nil, false, util.ChainError(err, "error creating crypt object")
		}
		break
	}

	w.printSuccessStatus("client authenticated")
	err := conn.SendSuccess()
	if err != nil {
		return nil, true, err
	}

	return c, false, nil
}
