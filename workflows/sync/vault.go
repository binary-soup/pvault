package syncworkflow

import (
	"github.com/binary-soup/go-commando/alert"
	"github.com/google/uuid"
)

const (
	UUID_SIZE = 16
)

type VaultItem struct {
	ID   uuid.UUID
	Name string
}

func NewVaultItem(id uuid.UUID, name string) VaultItem {
	return VaultItem{
		ID:   id,
		Name: name,
	}
}

func (v VaultItem) ToBytes() []byte {
	bytes := make([]byte, UUID_SIZE+len(v.Name))

	copy(bytes, v.ID[:])
	copy(bytes[UUID_SIZE:], []byte(v.Name))

	return bytes
}

func ParseVaultItemFromBytes(bytes []byte) (*VaultItem, error) {
	if len(bytes) <= UUID_SIZE {
		return nil, alert.Error("data too short")
	}

	id, err := uuid.FromBytes(bytes[:UUID_SIZE])
	if err != nil {
		return nil, alert.ChainError(err, "error loading uuid")
	}

	return &VaultItem{
		ID:   id,
		Name: string(bytes[UUID_SIZE:]),
	}, nil
}
