package vault

import (
	"strings"
)

type SearchItem struct {
	Name       string
	MatchStart int
	MatchEnd   int
}

func (v Vault) Search(substring string) []SearchItem {
	items := []SearchItem{}

	for name := range v.index {
		if substring == "" {
			items = append(items, v.newSearchItem(name, 0, 0))
		} else if idx := strings.Index(strings.ToLower(name), strings.ToLower(substring)); idx >= 0 {
			items = append(items, v.newSearchItem(name, idx, idx+len(substring)))
		}

	}

	return items
}

func (Vault) newSearchItem(name string, start, end int) SearchItem {
	return SearchItem{
		Name:       name,
		MatchStart: start,
		MatchEnd:   end,
	}
}
