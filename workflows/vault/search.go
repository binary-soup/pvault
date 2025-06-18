package vw

import (
	"fmt"
	"passwords/data/vault"
	"passwords/tools"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

var SEARCH_ITEM_STYLE = style.New(style.Yellow)
var SEARCH_ITEM_HIGHLIGHT = style.New(style.Yellow, style.Bold, style.Underline)

func (v VaultWorkflow) Search(search string) {
	for _, item := range v.Vault.Search(search) {
		v.styleSearchItem(item)
	}
}

func (v VaultWorkflow) SearchExactName(search string) (string, error) {
	items := v.Vault.Search(search)
	if len(items) == 0 {
		return "", util.Error(fmt.Sprintf("no matches found for \"%s\"", search))
	}

	if len(items) == 1 {
		fmt.Print("Match: ")
		v.styleSearchItem(items[0])

		return items[0].Name, nil
	}

	for idx, item := range items {
		style.Bolded.PrintF("[%d] ", idx+1)
		v.styleSearchItem(item)
	}

	n, err := tools.ReadInteger("INDEX", 1, len(items))
	if err != nil {
		return "", err
	}

	return items[n-1].Name, nil
}

func (v VaultWorkflow) styleSearchItem(item vault.SearchItem) {
	fmt.Printf("%s%s%s\n",
		SEARCH_ITEM_STYLE.Format(item.Name[:item.MatchStart]),
		SEARCH_ITEM_HIGHLIGHT.Format(item.Name[item.MatchStart:item.MatchEnd]),
		SEARCH_ITEM_STYLE.Format(item.Name[item.MatchEnd:]),
	)
}
