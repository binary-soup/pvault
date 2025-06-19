package sw

import (
	"fmt"
	"io"
	"net"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

const PORT = ":9000"

func (w SyncWorkflow) RunHost() error {
	// passkey, err := tools.ReadAndVerifyPasskey("Choose Host")
	// if err != nil {
	// 	return err
	// }

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

	bytes, err := io.ReadAll(conn)
	if err != nil {
		return util.ChainError(err, "error reading from connection")
	}

	fmt.Printf("%s: \"%s\"\n", style.Bolded.Format("MESSAGE"), string(bytes))

	return nil
}
