package config

import (
	"os"
	"path/filepath"
	"pvault/data/vault"

	"github.com/binary-soup/go-command/alert"
)

const PASSWORD_MIN_LIFETIME = 3

type Config struct {
	Vault    *vault.Vault   `json:"vault"`
	Passkey  PasskeyConfig  `json:"passkey"`
	Password PasswordConfig `json:"password"`
}

type PasskeyConfig struct {
	Timeout float32 `json:"timeout"`
}

type PasswordConfig struct {
	Lifetime float32 `json:"lifetime"`
}

func (cfg Config) Load() error {
	return cfg.Vault.Open()
}

func (cfg Config) Validate() ([]error, error) {
	errs := []error{}

	verr, err := cfg.validateVault()
	if err != nil {
		return nil, err
	}
	if verr != nil {
		errs = append(errs, verr)
	}

	if cfg.Passkey.Timeout < 0 {
		errs = append(errs, alert.Error("passkey timeout cannot be negative"))
	}

	if cfg.Password.Lifetime < PASSWORD_MIN_LIFETIME {
		errs = append(errs, alert.ErrorF("password lifetime cannot be less than %d seconds", PASSWORD_MIN_LIFETIME))
	}

	return errs, nil
}

func (cfg Config) validateVault() (error, error) {
	if cfg.Vault == nil {
		return alert.Error("missing vault config"), nil
	}

	if cfg.Vault.Path == "" {
		return alert.Error("vault path cannot be empty"), nil
	}

	_, err := os.Stat(filepath.Dir(cfg.Vault.Path))
	if os.IsNotExist(err) {
		return alert.Error("vault path does not exist"), nil
	}
	return nil, err
}
