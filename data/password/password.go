package password

import (
	"github.com/binary-soup/go-command/alert"
	"github.com/binary-soup/go-command/util"
)

type Password struct {
	Password      string   `json:"password,omitempty"`
	Username      string   `json:"username,omitempty"`
	URL           string   `json:"url,omitempty"`
	RecoveryCodes []string `json:"recovery_codes,omitempty"`
}

func LoadFile(path string) (*Password, error) {
	return util.LoadJSON[Password]("password", path)
}

func (password Password) SaveToFile(path string) error {
	return util.SaveJSON("password", &password, path)
}

func (password Password) Validate() error {
	if password.Password == "" && len(password.RecoveryCodes) == 0 {
		return alert.Error("both \"password\" and \"recovery codes\" cannot be empty")
	}
	return nil
}
