package sw

import (
	"fmt"
	"io"
	"net"
	"passwords/crypt"
	"passwords/tools"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

const PORT = ":9000"

func (w SyncWorkflow) RunHost() error {
	passkey, err := tools.ReadAndVerifyPasskey("Choose Host")
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		return util.ChainError(err, "error starting tcp server")
	}
	defer ln.Close()

	fmt.Printf("Listening on port %s\n", style.Bolded.Format(PORT))

	conn, err := ln.Accept()
	if err != nil {
		return util.ChainError(err, "error accepting client connection")
	}
	defer conn.Close()

	header := crypt.EmptyHeader()

	_, err = conn.Read(header)
	if err != nil {
		return util.ChainError(err, "error reading header from connection")
	}

	c, invalidPasskey, err := crypt.LoadCrypt(passkey, header)
	if invalidPasskey {
		//send ERROR response
		return util.Error("invalid client passkey")
	}
	if err != nil {
		return util.ChainError(err, "error creating crypt object")
	}

	ciphertext, err := io.ReadAll(conn)
	if err != nil {
		return util.ChainError(err, "error reading ciphertext from connection")
	}

	plaintext, err := c.Decrypt(ciphertext)
	if err != nil {
		return util.ChainError(err, "error decrypting ciphertext")
	}

	fmt.Printf("%s: \"%s\"\n", style.Bolded.Format("MESSAGE"), string(plaintext))

	//send SUCCUSS response
	return nil
}
