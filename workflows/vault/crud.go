package vw

import (
	"passwords/crypt"
	"passwords/data"
	"passwords/tools"

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

func (v VaultWorkflow) Encrypt(password *data.Password, passkey string) error {
	c, err := crypt.NewCrypt(passkey)
	if err != nil {
		return util.ChainError(err, "error initializing crypt tool")
	}

	bytes, err := password.Encrypt(c)
	if err != nil {
		return err
	}

	return v.Vault.SaveData(bytes, password.Name)
}

func (v VaultWorkflow) Decrypt(name string) (*data.Password, string, error) {
	bytes, err := v.Vault.ReadData(name)
	if err != nil {
		return nil, "", err
	}

	for {
		passkey, err := tools.ReadPasskey("Enter")
		if err != nil {
			return nil, "", err
		}

		c, invalidPasskey, err := crypt.LoadCrypt(passkey, bytes)
		if err != nil {
			return nil, "", err
		}
		if invalidPasskey {
			continue
		}

		password, err := data.DecryptPassword(c, bytes)
		if err != nil {
			return nil, "", err
		}

		return password, passkey, nil
	}
}

func (v VaultWorkflow) Delete(name string) error {
	return v.Vault.DeleteData(name)
}
