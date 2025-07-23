package vault

import (
	"github.com/binary-soup/go-commando/util"
)

type Cache struct {
	Passkey string `json:"passkey"`
}

func NewCache(passkey string) Cache {
	return Cache{
		Passkey: passkey,
	}
}

func (c Cache) SaveToFile(path string) error {
	return util.SaveJSON("cache", &c, path)
}

func LoadCacheFile(path string) (*Cache, error) {
	return util.LoadJSON[Cache]("cache", path)
}
