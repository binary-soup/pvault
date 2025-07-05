package syncworkflow

import (
	"pvault/crypt"
	"pvault/tools"
	"pvault/tools/sync"

	"github.com/binary-soup/go-command/alert"
)

func (w ClientWorkflow) authenticate(conn *sync.Connection) (*crypt.Crypt, error) {
	for {
		passkey, err := tools.ReadPasskey("Enter Host")
		if err != nil {
			return nil, err
		}

		crt, err := crypt.NewCrypt(passkey)
		if err != nil {
			return nil, alert.ChainError(err, "error creating crypt object")
		}

		conn.SendMessage("header", crt.Header)

		status, err := conn.ReadResponse()
		if status == sync.SUCCESS {
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

func (w HostWorkflow) authenticate(conn *sync.Connection, passkey string) (*crypt.Crypt, error) {
	for {
		header, err := conn.ReadMessage("header")
		if err != nil {
			conn.SendClientError("error reading crypt header message")
			return nil, err
		}

		crt, invalidPasskey, err := crypt.LoadCrypt(passkey, header)
		if invalidPasskey {
			errorLog.Log("invalid client passkey")
			conn.SendAuthError("invalid passkey")
			continue
		}
		if err != nil {
			conn.SendInternalError()
			return nil, alert.ChainError(err, "error creating crypt object")
		}

		conn.SendSuccess()
		successLog.Log("client authenticated")

		return crt, nil
	}
}
