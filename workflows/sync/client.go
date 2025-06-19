package sw

import (
	"fmt"
	"net"

	"github.com/binary-soup/go-command/util"
)

func (w SyncWorkflow) RunClient(addr string) error {
	// passkey, err := tools.ReadPasskey("Enter Host")
	// if err != nil {
	// 	return err
	// }

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return util.ChainErrorF(err, "error dialing host at \"%s\"", addr)
	}
	defer conn.Close()

	fmt.Fprint(conn, "This is a test")
	return nil
}
