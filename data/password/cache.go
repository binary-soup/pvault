package password

import "github.com/binary-soup/go-command/util"

type Cache struct {
	Password *Password `json:"password"`
	Meta     *Meta     `json:"meta"`
}

func LoadCacheFile(path string) (*Cache, error) {
	return util.LoadJSON[Cache]("password cache", path)
}

func (cache Cache) SaveToFile(path string) error {
	return util.SaveJSON("cache", &cache, path)
}
