package workflows

import (
	"passwords/crypt"
	"passwords/data"
	"passwords/tools"

	"github.com/binary-soup/go-command/util"
)

func EncryptToVault(vault data.Vault, password *data.Password, filename string) error {
	var err error

	if password.Passkey == "" {
		password.Passkey, err = tools.ReadAndVerifyPasskey("Choose New")
	} else {
		err = tools.VerifyPasskey(password.Passkey)
	}

	if err != nil {
		return err
	}

	c, err := crypt.NewCrypt(password.Passkey)
	if err != nil {
		return util.ChainError(err, "error initializing crypt tool")
	}

	bytes, err := password.Encrypt(c)
	if err != nil {
		return err
	}

	err = vault.SaveData(bytes, filename+".crypt")
	if err != nil {
		return err
	}

	return nil
}

func DecryptFromVault(vault data.Vault, filename string) (*data.Password, error) {
	bytes, err := vault.LoadData(filename + ".crypt")
	if err != nil {
		return nil, err
	}

	for {
		passkey, err := tools.ReadPasskey("Enter")
		if err != nil {
			return nil, err
		}

		c, invalidPasskey, err := crypt.LoadCrypt(passkey, bytes)
		if err != nil {
			return nil, err
		}
		if invalidPasskey {
			continue
		}

		return data.DecryptPassword(c, bytes)
	}
}

func DeleteFromVault(vault data.Vault, filename string) error {
	return vault.Delete(filename + ".crypt")
}

func SearchVault(vault data.Vault, search string) ([]string, error) {
	return vault.Search(search, ".crypt")
}
