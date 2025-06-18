package workflows

import (
	"fmt"
	"passwords/crypt"
	"passwords/data"
	"passwords/data/vault"
	"passwords/tools"

	"github.com/binary-soup/go-command/util"
)

func vaultFilename(id uint) string {
	return fmt.Sprintf("u%d.crypt", id)
}

func EncryptToVault(v vault.Vault, password *data.Password, index *vault.Index) error {
	var err error

	if index.Passkey == "" {
		index.Passkey, err = tools.ReadAndVerifyPasskey("Choose New")
	} else {
		err = tools.VerifyPasskey(index.Passkey)
	}

	if err != nil {
		return err
	}

	c, err := crypt.NewCrypt(index.Passkey)
	if err != nil {
		return util.ChainError(err, "error initializing crypt tool")
	}

	bytes, err := password.Encrypt(c)
	if err != nil {
		return err
	}

	err = v.SaveData(bytes, vaultFilename(index.ID))
	if err != nil {
		return err
	}

	return nil
}

func DecryptFromVault(v vault.Vault, id uint) (*data.Password, *vault.Index, error) {
	bytes, err := v.LoadData(vaultFilename(id))
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

		index := &vault.Index{
			ID:      id,
			Passkey: passkey,
		}

		return password, index, nil
	}
}

func DeleteFromVault(vault vault.Vault, id uint) error {
	return vault.Delete(vaultFilename(id))
}
