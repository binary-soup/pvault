package sw

import (
	"net"
	"passwords/crypt"
	"passwords/tools"

	"github.com/binary-soup/go-command/util"
)

func (w SyncWorkflow) RunClient(addr string) error {
	const MESSAGE = "This is a test message."

	passkey, err := tools.ReadPasskey("Enter Host")
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return util.ChainErrorF(err, "error dialing host at \"%s\"", addr)
	}
	defer conn.Close()

	c, err := crypt.NewCrypt(passkey)
	if err != nil {
		return util.ChainError(err, "error creating crypt object")
	}

	ciphertext := c.Encrypt([]byte(MESSAGE))

	_, err = conn.Write(c.Header)
	if err != nil {
		return util.ChainError(err, "error writing header to connection")
	}

	_, err = conn.Write(ciphertext)
	if err != nil {
		return util.ChainError(err, "error writing ciphertext to connection")
	}

	return nil
}
