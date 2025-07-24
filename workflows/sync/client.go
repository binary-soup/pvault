package syncworkflow

import (
	"fmt"
	"pvault/data/vault"
	"pvault/tools"
	"pvault/tools/sync"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/style"
)

type ClientWorkflow struct {
	Vault *vault.Vault
}

func NewClientWorkflow(v *vault.Vault) ClientWorkflow {
	return ClientWorkflow{
		Vault: v,
	}
}

func (w ClientWorkflow) Run(addr, port string) error {
	client := sync.NewClient()

	conn, err := client.Connect(addr + port)
	if err != nil {
		return err
	}
	defer conn.Close()

	return w.accept(conn)
}

func (w ClientWorkflow) accept(conn *sync.Connection) error {
	fmt.Printf("Connected to %s\n", style.BoldInfo.Format(conn.RemoteAddress()))

	hostname, err := conn.ExchangeHostname()
	if err != nil {
		return err
	}
	successLog.LogF("host identified as %s", style.BoldInfo.Format(hostname))

	conn.Crypt, err = w.authenticate(conn)
	if err != nil {
		return err
	}

	successLog.Log("receiving vault list")
	list, err := conn.ReadMessageBlock("vault list")
	if err != nil {
		return err
	}

	for i, bytes := range list {
		style.Bolded.PrintF("[%d/%d] ", i+1, len(list))

		err := w.requestVaultItem(conn, bytes)
		if err != nil {
			errorLog.Log(err)
		}
		w.Vault.Flush()
	}
	conn.SendNoRequest()

	successLog.Log("end of list")
	return nil
}

func (w ClientWorkflow) hostError(err error) error {
	return alert.ChainError(err, "error from host")
}

func (w ClientWorkflow) requestVaultItem(conn *sync.Connection, bytes []byte) error {
	item, err := ParseVaultItemFromBytes(bytes)
	if err != nil {
		return alert.ChainError(err, "error parsing vault item from vault list")
	}

	style.Highlight.Println(item.Name)

	if !w.promptAcceptItem(item) {
		errorLog.LogF("denied item %s", style.Bolded.Format(item.ID.String()))
		return nil
	}

	successLog.LogF("accepted item %s", style.Bolded.Format(item.ID.String()))
	conn.SendRequest("vault item", 1, item.ID[:])

	status, err := conn.ReadResponse()
	if status != sync.SUCCESS {
		return w.hostError(err)
	}
	if err != nil {
		return err
	}

	bytes, err = conn.ReadMessage("vault item")
	if err != nil {
		return err
	}

	err = w.Vault.SaveRaw(item.ID, item.Name, bytes)
	if err != nil {
		return err
	}

	successLog.Log("item received")
	return nil
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
	style.Info.Println("(up to date)")
	return false
}
