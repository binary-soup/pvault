package vw

import "pvault/data/vault"

type VaultWorkflow struct {
	Vault vault.Vault
}

func NewVaultWorkflow(v vault.Vault) VaultWorkflow {
	return VaultWorkflow{
		Vault: v,
	}
}
