package vault

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/binary-soup/go-commando/alert"
	"github.com/google/uuid"
)

const INDEX_FILE = "index.txt"

type Index struct {
	nameMap map[string]uuid.UUID
	uuidMap map[uuid.UUID]string

	isDirty bool
}

func newIndex() *Index {
	return &Index{
		nameMap: map[string]uuid.UUID{},
		uuidMap: map[uuid.UUID]string{},
		isDirty: false,
	}
}

func (idx Index) Count() int {
	return len(idx.nameMap)
}

func (idx *Index) AddPair(name string, id uuid.UUID) {
	oldName, ok := idx.uuidMap[id]
	if ok && oldName != name {
		delete(idx.nameMap, oldName)
	}

	idx.nameMap[name] = id
	idx.uuidMap[id] = name

	idx.isDirty = true
}

func (idx Index) HasName(name string) bool {
	_, ok := idx.nameMap[name]
	return ok
}

func (idx Index) HasID(id uuid.UUID) bool {
	_, ok := idx.uuidMap[id]
	return ok
}

func (idx Index) GetName(id uuid.UUID) (string, error) {
	name, ok := idx.uuidMap[id]
	if !ok {
		return "", alert.ErrorF("id \"%s\" not found", id.String())
	}
	return name, nil
}

func (idx Index) GetID(name string) (uuid.UUID, error) {
	id, ok := idx.nameMap[name]
	if !ok {
		return uuid.Nil, alert.ErrorF("name \"%s\" not found", name)
	}
	return id, nil
}

func (idx *Index) DeleteID(id uuid.UUID) {
	name := idx.uuidMap[id]

	delete(idx.uuidMap, id)
	delete(idx.nameMap, name)

	idx.isDirty = true
}

func (idx Index) Iterate(itr func(int, string, uuid.UUID)) {
	i := 0
	for name, id := range idx.nameMap {
		itr(i, name, id)
		i++
	}
}

func (v Vault) saveIndex() error {
	if !v.Index.isDirty {
		return nil
	}

	file, err := os.Create(filepath.Join(v.Path, INDEX_FILE))
	if err != nil {
		return alert.ChainError(err, "error creating index file")
	}
	defer file.Close()

	v.Index.Iterate(func(_ int, name string, id uuid.UUID) {
		fmt.Fprintf(file, "%s:%s\n", id.String(), name)
	})

	v.Index.isDirty = false
	return nil
}

func (v Vault) loadIndex() (*Index, error) {
	file, err := os.Open(filepath.Join(v.Path, INDEX_FILE))
	if os.IsNotExist(err) {
		return newIndex(), nil
	}
	if err != nil {
		return nil, alert.ChainError(err, "error opening index file")
	}
	defer file.Close()

	index := newIndex()
	line := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line++

		name, id, err := v.parseIndexPair(scanner.Text(), line)
		if err != nil {
			return nil, err
		}
		index.AddPair(name, id)
	}

	if err := scanner.Err(); err != nil {
		return nil, alert.ChainError(err, "error parsing index file")
	}
	return index, nil
}

func (v Vault) parseIndexPair(line string, lineNumber int) (string, uuid.UUID, error) {
	tokens := strings.SplitN(line, ":", 2)
	if len(tokens) < 2 {
		return "", uuid.Nil, alert.ErrorF("[line %d] invalid index pair", lineNumber)
	}

	id, err := uuid.Parse(tokens[0])
	if err != nil {
		return "", uuid.Nil, alert.ChainErrorF(err, "[line %d] invalid uuid", lineNumber)
	}

	return tokens[1], id, nil
}
