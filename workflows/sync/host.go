package syncworkflow

import (
	"fmt"
	"pvault/data/vault"
	"pvault/tools"
	"pvault/tools/sync"

	"github.com/binary-soup/go-command/alert"
	"github.com/binary-soup/go-command/style"
	"github.com/google/uuid"
)

type HostWorkflow struct {
	Vault vault.Vault
}

func NewHostWorkflow(v vault.Vault) HostWorkflow {
	return HostWorkflow{
		Vault: v,
	}
}

func (w HostWorkflow) Run(port string, persist bool) error {
	passkey, err := tools.ReadAndVerifyPasskey("Choose Host")
	if err != nil {
		return err
	}

	host := sync.NewHost(port)

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

		fmt.Printf("Connected with %s\n", style.BoldInfo.Format(conn.RemoteAddress()))

		err = w.accept(conn, passkey)
		if !persist {
			return err
		}
		if err != nil {
			errorLog.Log(err)
		}
	}
}

func (w HostWorkflow) accept(conn *sync.Connection, passkey string) error {
	hostname, err := conn.ExchangeHostname()
	if err != nil {
		return err
	}
	successLog.LogF("client identified as %s", style.BoldInfo.Format(string(hostname)))

	conn.Crypt, err = w.authenticate(conn, passkey)
	if err != nil {
		return err
	}

	list := w.buildVaultList()
	conn.SendMessageBlock("vault list", list)
	successLog.LogF("sent vault list (%d)", len(list))

	err = conn.ReceiveManyRequests(func(_ int, msg []byte) {
		err := w.sendVaultFile(conn, msg)
		if err != nil {
			errorLog.Log(err)
		}
	})
	if err != nil {
		return err
	}

	successLog.Log("end of requests")
	return nil
}

func (w HostWorkflow) buildVaultList() [][]byte {
	items := make([][]byte, w.Vault.Index.Count())

	w.Vault.Index.Iterate(func(i int, name string, id uuid.UUID) {
		items[i] = NewVaultItem(id, name).ToBytes()
	})

	return items
}

func (w HostWorkflow) sendVaultFile(conn *sync.Connection, bytes []byte) error {
	id, err := uuid.FromBytes(bytes)
	if err != nil {
		conn.SendClientError("could not parse uuid")
		return alert.ChainError(err, "error parsing message uuid")
	}

	if !w.Vault.Index.HasID(id) {
		conn.SendClientError("invalid uuid")
		return alert.ErrorF("received invalid uuid %s", id.String())
	}

	successLog.LogF("received vault request for %s", style.Bolded.Format(id.String()))

	raw, err := w.Vault.LoadRaw(id)
	if err != nil {
		conn.SendInternalError()
		return err
	}

	conn.SendSuccess()
	conn.SendMessage("vault item", raw)

	successLog.LogF("item sent (%d bytes)", len(raw))
	return nil
}
