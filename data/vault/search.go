package vault

import (
	"slices"
	"strings"

	"github.com/google/uuid"
)

type SearchItem struct {
	Name       string
	MatchStart int
	MatchEnd   int
}

func (v Vault) Search(substring string) []SearchItem {
	items := []SearchItem{}

	v.Index.Iterate(func(_ int, name string, _ uuid.UUID) {
		if substring == "" {
			items = append(items, v.newSearchItem(name, 0, 0))
		} else if idx := strings.Index(strings.ToLower(name), strings.ToLower(substring)); idx >= 0 {
			items = append(items, v.newSearchItem(name, idx, idx+len(substring)))
		}
	})

	slices.SortFunc(items, func(a, b SearchItem) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	return items
}

func (Vault) newSearchItem(name string, start, end int) SearchItem {
	return SearchItem{
		Name:       name,
		MatchStart: start,
		MatchEnd:   end,
	}
}
