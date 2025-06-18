package vault

import "github.com/binary-soup/go-command/util"

type Index struct {
	ID      uint   `json:"id"`
	Passkey string `json:"passkey"`
}

func (idx Index) SaveToFile(path string) error {
	return util.SaveJSON("index", &idx, path)
}

func LoadIndexFile(path string) (*Index, error) {
	return util.LoadJSON[Index]("index", path)
}
