package vault

import (
	"strings"
)

type SearchItem struct {
	Name       string
	ID         uint
	MatchStart int
	MatchEnd   int
}

func (v Vault) Search(substring string) []SearchItem {
	items := []SearchItem{}

	for key, val := range v.indexMap {
		if substring == "" {
			items = append(items, v.newSearchItem(key, val, 0, 0))
		} else if idx := strings.Index(strings.ToLower(key), strings.ToLower(substring)); idx >= 0 {
			items = append(items, v.newSearchItem(key, val, idx, idx+len(substring)))
		}

	}

	return items
}

func (Vault) newSearchItem(name string, id uint, start, end int) SearchItem {
	return SearchItem{
		Name:       name,
		ID:         id,
		MatchStart: start,
		MatchEnd:   end,
	}
}
