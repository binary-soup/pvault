package workflows

import (
	"passwords/crypt"
	"passwords/data"
	"passwords/tools"

	"github.com/binary-soup/go-command/util"
)

func EncryptToVault(vault data.Vault, password *data.Password, index *data.Index, filename string) error {
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

	err = vault.SaveData(bytes, filename+".json.crypt")
	if err != nil {
		return err
	}

	return nil
}

func DecryptFromVault(vault data.Vault, filename string) (*data.Password, *data.Index, error) {
	bytes, err := vault.LoadData(filename + ".json.crypt")
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
			Passkey: passkey,
		}

		return password, index, nil
	}
}

func DeleteFromVault(vault data.Vault, filename string) error {
	return vault.Delete(filename + ".crypt")
}

func SearchVault(vault data.Vault, search string) ([]string, error) {
	return vault.Search(search, ".crypt")
}
