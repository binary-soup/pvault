package workflows

import (
	"fmt"
	"passwords/crypt"
	"passwords/data"
	"passwords/tools"

	"github.com/binary-soup/go-command/util"
)

func vaultFilename(id uint) string {
	return fmt.Sprintf("u%d.crypt", id)
}

func EncryptToVault(vault data.Vault, password *data.Password, index *data.Index) error {
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

	err = vault.SaveData(bytes, vaultFilename(index.ID))
	if err != nil {
		return err
	}

	return nil
}

func DecryptFromVault(vault data.Vault, id uint) (*data.Password, *data.Index, error) {
	bytes, err := vault.LoadData(vaultFilename(id))
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

		index := &data.Index{
			ID:      id,
			Passkey: passkey,
		}

		return password, index, nil
	}
}

func DeleteFromVault(vault data.Vault, id uint) error {
	return vault.Delete(vaultFilename(id))
}

func SearchVault(vault data.Vault, search string) ([]string, error) {
	return vault.Search(search, ".crypt")
}
