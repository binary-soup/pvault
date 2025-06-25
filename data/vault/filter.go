package vault

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
	"github.com/google/uuid"
)

const FILTER_FILE = "filter.txt"

type Filter struct {
	uuidSet map[uuid.UUID]struct{}
	isDirty bool
}

func newFilter() *Filter {
	return &Filter{
		uuidSet: map[uuid.UUID]struct{}{},
		isDirty: false,
	}
}

func (f *Filter) AddItem(id uuid.UUID) {
	f.uuidSet[id] = struct{}{}
	f.isDirty = true
}

func (f Filter) IsFiltered(id uuid.UUID) bool {
	_, ok := f.uuidSet[id]
	return ok
}

func (f *Filter) Clear() {
	f.uuidSet = map[uuid.UUID]struct{}{}
	f.isDirty = true
}

func (f Filter) Iterate(itr func(int, uuid.UUID)) {
	i := 0
	for id := range f.uuidSet {
		itr(i, id)
		i++
	}
}

func (v Vault) saveFilter() error {
	if !v.Filter.isDirty {
		return nil
	}

	file, err := os.Create(filepath.Join(v.Path, FILTER_FILE))
	if err != nil {
		return util.ChainError(err, "error creating filter file")
	}
	defer file.Close()

	v.Filter.Iterate(func(_ int, id uuid.UUID) {
		fmt.Fprintln(file, id.String())
	})

	v.Filter.isDirty = false
	return nil
}

func (v Vault) loadFilter() (*Filter, error) {
	file, err := os.Open(filepath.Join(v.Path, FILTER_FILE))
	if os.IsNotExist(err) {
		return newFilter(), nil
	}
	if err != nil {
		return nil, util.ChainError(err, "error opening filter file")
	}
	defer file.Close()

	filter := newFilter()
	line := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line++

		id, err := uuid.Parse(scanner.Text())
		if err != nil {
			return nil, util.ChainErrorF(err, "[line %d] invalid uuid", line)
		}
		filter.AddItem(id)
	}

	if err := scanner.Err(); err != nil {
		return nil, util.ChainError(err, "error parsing filter file")
	}
	return filter, nil
}
