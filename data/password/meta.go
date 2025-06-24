package password

import "github.com/google/uuid"

type Meta struct {
	Name    string    `json:"name"`
	Passkey string    `json:"passkey"`
	ID      uuid.UUID `json:"uuid"`
}

func NewMeta(name, passkey string) *Meta {
	return &Meta{
		Name:    name,
		Passkey: passkey,
		ID:      uuid.New(),
	}
}
