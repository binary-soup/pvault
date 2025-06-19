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

	conn, err := host.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	header, err := conn.ReadMessage("header")
	if err != nil {
		return err
	}

	c, invalidPasskey, err := crypt.LoadCrypt(passkey, header)
	if invalidPasskey {
		//send ERROR response
		return util.Error("invalid client passkey")
	}
	if err != nil {
		return util.ChainError(err, "error creating crypt object")
	}

	ciphertext, err := conn.ReadMessage("ciphertext")
	if err != nil {
		return err
	}

	plaintext, err := c.Decrypt(ciphertext)
	if err != nil {
		return util.ChainError(err, "error decrypting ciphertext")
	}

	fmt.Printf("%s: \"%s\"\n", style.Bolded.Format("MESSAGE"), string(plaintext))

	//send SUCCESS response
	return nil
}
