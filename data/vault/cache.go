package vault

import (
	"github.com/binary-soup/go-command/util"
)

type Cache struct {
	Passkey string `json:"passkey"`
}

func (c Cache) SaveToFile(path string) error {
	return util.SaveJSON("index", &c, path)
}

func LoadCacheFile(path string) (*Cache, error) {
	return util.LoadJSON[Cache]("cache", path)
}
