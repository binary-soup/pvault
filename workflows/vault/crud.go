package vw

import (
	"fmt"
	"passwords/crypt"
	"passwords/data"
	"passwords/data/vault"
	"passwords/tools"

	"github.com/binary-soup/go-command/util"
	"github.com/google/uuid"
)

func (v VaultWorkflow) Encrypt(password *data.Password, cache *vault.Cache) error {
	if v.Vault.NameExists(password.Name) {
		return util.Error(fmt.Sprintf("name \"%s\" already exists", password.Name))
	}

	var err error
	if cache.Passkey == "" {
		cache.Passkey, err = tools.ReadAndVerifyPasskey("Choose New")
	} else {
		err = tools.VerifyPasskey(cache.Passkey)
	}

	if err != nil {
		return err
	}

	c, err := crypt.NewCrypt(cache.Passkey)
	if err != nil {
		return util.ChainError(err, "error initializing crypt tool")
	}

	bytes, err := password.Encrypt(c)
	if err != nil {
		return err
	}

	return v.Vault.SaveData(bytes, uuid.New(), password.Name)
}

func (v VaultWorkflow) Decrypt(name string) (*data.Password, *vault.Cache, error) {
	bytes, err := v.Vault.ReadData(name)
	if err != nil {
		return nil, nil, err
	}

	for {
		passkey, err := tools.ReadPasskey("Enter")
		if err != nil {
			return nil, nil, err
		}

		c, invalidPasskey, err := crypt.LoadCrypt(passkey, bytes)
		if err != nil {
			return nil, nil, err
		}
		if invalidPasskey {
			continue
		}

		password, err := data.DecryptPassword(c, bytes)
		if err != nil {
			return nil, nil, err
		}

		index := &vault.Cache{
			Passkey: passkey,
		}

		return password, index, nil
	}
}

func (v VaultWorkflow) Delete(name string) error {
	return v.Vault.DeleteData(name)
}
