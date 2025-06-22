package data

import (
	"os"
	"path/filepath"
	"pvault/data/vault"

	"github.com/binary-soup/go-command/util"
)

const CONFIG_PATH = "config.json"

type Config struct {
	Vault   vault.Vault   `json:"vault"`
	Passkey PasskeyConfig `json:"passkey"`
}

type PasskeyConfig struct {
	Timeout float32 `json:"timeout"`
}

func LoadBaseConfig() (*Config, error) {
	path, _ := os.Executable()
	path = filepath.Dir(path)

	return LoadCustomConfig(filepath.Join(path, CONFIG_PATH))
}

func LoadCustomConfig(path string) (*Config, error) {
	cfg, err := util.LoadJSON[Config]("config", path)
	if err != nil {
		return nil, err
	}

	err = cfg.Vault.Open()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
