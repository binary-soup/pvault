package config

import (
	"os"
	"path/filepath"
	"pvault/data/vault"

	"github.com/binary-soup/go-command/alert"
)

type Config struct {
	Vault   *vault.Vault  `json:"vault"`
	Passkey PasskeyConfig `json:"passkey"`
}

type PasskeyConfig struct {
	Timeout float32 `json:"timeout"`
}

func (cfg Config) Load() error {
	return cfg.Vault.Open()
}

func (cfg Config) Validate() ([]error, error) {
	errs := []error{}

	if cfg.Vault.Path == "" {
		errs = append(errs, alert.Error("vault path cannot be empty"))
	} else {
		_, err := os.Stat(filepath.Dir(cfg.Vault.Path))
		if os.IsNotExist(err) {
			errs = append(errs, alert.Error("vault path does not exist"))
		} else if err != nil {
			return nil, err
		}
	}

	if cfg.Passkey.Timeout < 0 {
		errs = append(errs, alert.Error("passkey timeout cannot be negative"))
	}

	return errs, nil
}
