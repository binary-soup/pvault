package data

import (
	"os"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
)

type Config struct {
	Vault Vault `json:"vault"`
}

const CONFIG_PATH = "config.json"

func LoadConfig() (*Config, error) {
	path, _ := os.Executable()
	path = filepath.Dir(path)

	cfg, err := util.LoadJSON[Config]("config", filepath.Join(path, CONFIG_PATH))
	if err != nil {
		return nil, err
	}

	err = cfg.Vault.Init()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
