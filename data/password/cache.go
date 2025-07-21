package password

import "github.com/binary-soup/go-commando/data"

type Cache struct {
	Password Password `json:"password"`
	Meta     Meta     `json:"meta"`
}

func LoadCacheFile(path string) (Cache, error) {
	return data.LoadJSON[Cache]("password cache", path)
}

func (cache Cache) SaveToFile(path string) error {
	return data.SaveJSON("cache", &cache, path)
}
