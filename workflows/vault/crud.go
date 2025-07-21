package vw

import (
	"pvault/crypt"
	"pvault/data/password"
	"pvault/tools"

	"github.com/binary-soup/go-commando/alert"
)

func (v VaultWorkflow) Encrypt(cache password.Cache) error {
	c, err := crypt.NewCrypt(cache.Meta.Passkey)
	if err != nil {
		return alert.ChainError(err, "error initializing crypt tool")
	}

	ciphertext, err := cache.Password.Encrypt(c)
	if err != nil {
		return err
	}

	return v.Vault.SaveData(c.Header, ciphertext, cache.Meta.ID, cache.Meta.Name)
}

func (v VaultWorkflow) Decrypt(name string, timeout float32) (*password.Cache, error) {
	id, err := v.Vault.Index.GetID(name)
	if err != nil {
		return nil, err
	}

	header, ciphertext, err := v.Vault.LoadData(id)
	if err != nil {
		return nil, err
	}

	for {
		passkey, err := tools.ReadPasskey("Enter")
		if err != nil {
			return nil, err
		}

		c, invalidPasskey, err := crypt.LoadCrypt(passkey, header)
		if err != nil {
			return nil, err
		}
		if invalidPasskey {
			tools.Timeout(timeout)
			continue
		}

		pswrd, err := password.Decrypt(c, ciphertext)
		if err != nil {
			return nil, err
		}

		meta := password.Meta{
			Name:    name,
			Passkey: passkey,
			ID:      id,
		}

		return &password.Cache{
			Password: pswrd,
			Meta:     meta,
		}, nil
	}
}

func (v VaultWorkflow) Delete(name string) error {
	id, err := v.Vault.Index.GetID(name)
	if err != nil {
		return err
	}

	return v.Vault.DeleteData(id)
}
