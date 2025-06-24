package syncworkflow

import (
	"fmt"
	"pvault/crypt"
	"pvault/data/vault"
	"pvault/tools"
	"pvault/tools/sync"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type ClientWorkflow struct {
	Vault vault.Vault
}

func NewClientWorkflow(v vault.Vault) ClientWorkflow {
	return ClientWorkflow{
		Vault: v,
	}
}

func (w ClientWorkflow) Run(addr, port string) error {
	client := sync.NewClient()

	conn, err := client.Connect(addr + port)
	if err != nil {
		return nil
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", style.BoldInfo.Format(conn.RemoteAddress()))
	return w.accept(conn)
}

func (w ClientWorkflow) accept(conn *sync.Connection) error {
	hostname, err := conn.ExchangeHostname()
	if err != nil {
		return err
	}
	successLog.LogF("host identified as %s", style.BoldInfo.Format(hostname))

	crt, err := w.authenticate(conn)
	if err != nil {
		return err
	}

	successLog.Log("receiving vault list")
	err = conn.ReadManyMessages("vault list", crt, func(idx, count uint32, bytes []byte) error {
		item, err := ParseVaultItemFromBytes(bytes)
		if err != nil {
			return util.ChainError(err, "error parsing vault item from vault list")
		}

		fmt.Printf("%s %s ", style.Bolded.FormatF("[%d/%d]", idx+1, count), NAME_STYLE.Format(item.Name))

		if w.promptAcceptItem(item) {
			successLog.LogF("accepted item %s", item.ID.String())
		} else {
			errorLog.LogF("denied item %s", item.ID.String())
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (w ClientWorkflow) hostError(err error) error {
	return util.ChainError(err, "error from host")
}

func (w ClientWorkflow) authenticate(conn *sync.Connection) (*crypt.Crypt, error) {
	for {
		passkey, err := tools.ReadPasskey("Enter Host")
		if err != nil {
			return nil, err
		}

		crt, err := crypt.NewCrypt(passkey)
		if err != nil {
			return nil, util.ChainError(err, "error creating crypt object")
		}

		conn.SendMessage("header", crt.Header)

		status, err := conn.ReadResponse()
		if status == sync.ERROR_NONE {
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

func (w ClientWorkflow) promptAcceptItem(item *VaultItem) bool {
	if w.Vault.Filter.IsFiltered(item.ID) {
		style.Info.Println("(filtered)")
		return false
	}

	if w.Vault.Index.HasID(item.ID) {
		return w.promptUpdateItem(item)
	} else {
		return w.promptNewItem(item)
	}
}

func (w ClientWorkflow) promptNewItem(item *VaultItem) bool {
	res := tools.PromptAccept(fmt.Sprintf("%s [y/n/N]?", []byte(style.Create.Format("(accept new file)"))), []byte("ynN"))
	if res == 1 {
		return false
	}
	if res == 2 {
		w.Vault.Filter.AddItem(item.ID)
		style.Info.Println("(item added to filter)")
		return false
	}

	for w.Vault.Index.HasName(item.Name) {
		item.Name = tools.PromptString(false, fmt.Sprintf("%s Choose new %s:", style.Info.Format("(name in use)"), style.Bolded.Format("NAME")))
	}

	return true
}

func (w ClientWorkflow) promptUpdateItem(item *VaultItem) bool {
	//TODO: implement modified time
	style.Info.PrintF("%s up to date\n", NAME_STYLE.Format(item.Name))
	return false
}
