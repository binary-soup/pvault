package vault

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/binary-soup/go-command/util"
	"github.com/google/uuid"
)

type indexMap map[string]uuid.UUID

type Vault struct {
	Path string `json:"path"`

	index indexMap
}

func (v *Vault) Init() error {
	err := os.Mkdir(v.Path, 0755)
	if err != nil && !os.IsExist(err) {
		return util.ChainError(err, "error creating vault directory")
	}

	v.index, err = v.loadIndex()
	if err != nil {
		return util.ChainError(err, "error loading index map")
	}

	return nil
}

func (v Vault) loadIndex() (indexMap, error) {
	file, err := os.Open(filepath.Join(v.Path, INDEX_FILE))
	if os.IsNotExist(err) {
		return indexMap{}, nil
	}
	if err != nil {
		return nil, util.ChainError(err, "error opening index file")
	}
	defer file.Close()

	index := indexMap{}

	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		line++

		name, id, err := v.parseIndexPair(scanner.Text(), line)
		if err != nil {
			return nil, err
		}
		index[name] = id
	}

	if err := scanner.Err(); err != nil {
		return nil, util.ChainError(err, "error parsing index file")
	}
	return index, nil
}

func (v Vault) parseIndexPair(line string, lineNumber int) (string, uuid.UUID, error) {
	tokens := strings.SplitN(line, ":", 2)
	if len(tokens) < 2 {
		return "", uuid.Nil, util.Error(fmt.Sprintf("[line %d] invalid index pair", lineNumber))
	}

	id, err := uuid.Parse(tokens[0])
	if err != nil {
		return "", uuid.Nil, util.ChainErrorF(err, "[line %d] invalid uuid", lineNumber)
	}

	return tokens[1], id, nil
}
