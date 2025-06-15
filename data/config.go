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

	err = os.Mkdir(cfg.Vault.Path, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, util.ChainError(err, "error creating vault directory")
	}

	return cfg, nil
}
