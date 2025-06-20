package vw

import (
	"pvault/crypt"
	"pvault/data"
	"pvault/tools"

	"github.com/binary-soup/go-command/util"
)

func (v VaultWorkflow) ChooseOrVerifyPasskey(passkey *string) error {
	var err error

	if *passkey == "" {
		*passkey, err = tools.ReadAndVerifyPasskey("Choose New")
	} else {
		err = tools.VerifyPasskey(*passkey)
	}

	return err
}

func (v VaultWorkflow) Encrypt(password *data.Password, cache *data.PasswordCache) error {
	c, err := crypt.NewCrypt(cache.Passkey)
	if err != nil {
		return util.ChainError(err, "error initializing crypt tool")
	}

	ciphertext, err := password.Encrypt(c)
	if err != nil {
		return err
	}

	return v.Vault.SaveData(c.Header, ciphertext, cache.ID, password.Name)
}

func (v VaultWorkflow) Decrypt(name string, timeout float32) (*data.Password, *data.PasswordCache, error) {
	header, ciphertext, id, err := v.Vault.ReadData(name)
	if err != nil {
		return nil, nil, err
	}

	for {
		passkey, err := tools.ReadPasskey("Enter")
		if err != nil {
			return nil, nil, err
		}

		c, invalidPasskey, err := crypt.LoadCrypt(passkey, header)
		if err != nil {
			return nil, nil, err
		}
		if invalidPasskey {
			tools.Timeout(timeout)
			continue
		}

		password, err := data.DecryptPassword(c, ciphertext)
		if err != nil {
			return nil, nil, err
		}

		cache := &data.PasswordCache{
			Passkey: passkey,
			ID:      id,
		}

		return password, cache, nil
	}
}

func (v VaultWorkflow) Delete(name string) error {
	return v.Vault.DeleteData(name)
}
