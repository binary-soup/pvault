package vw

import (
	"fmt"
	"pvault/data/vault"
	"pvault/tools"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/style"
)

var SEARCH_ITEM_STYLE = style.New(style.Yellow)
var SEARCH_ITEM_HIGHLIGHT = style.New(style.Yellow, style.Bold, style.Underline)

func (v VaultWorkflow) Search(search string) []vault.SearchItem {
	items := v.Vault.Search(search)
	pad := len(fmt.Sprintf("%d", len(items)))

	for idx, item := range items {
		v.styleSearchItem(item, idx+1, pad)
	}
	return items
}

func (v VaultWorkflow) SearchExactName(search string) (string, error) {
	items := v.Search(search)
	if len(items) == 0 {
		return "", alert.ErrorF("no matches found for \"%s\"", search)
	}

	if len(items) == 1 {
		return items[0].Name, nil
	}

	n, err := tools.ReadInteger("INDEX", 1, len(items))
	if err != nil {
		return "", err
	}

	return items[n-1].Name, nil
}

func (v VaultWorkflow) styleSearchItem(item vault.SearchItem, idx, pad int) {
	fmt.Printf("%s %s%s%s\n",
		style.Bolded.FormatF("[%*d]", pad, idx),
		SEARCH_ITEM_STYLE.Format(item.Name[:item.MatchStart]),
		SEARCH_ITEM_HIGHLIGHT.Format(item.Name[item.MatchStart:item.MatchEnd]),
		SEARCH_ITEM_STYLE.Format(item.Name[item.MatchEnd:]),
	)
}
