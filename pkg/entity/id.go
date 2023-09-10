package entity

import (
	"github.com/google/uuid"
)

type ID = uuid.UUID

func NewID() ID {
	return uuid.New()
}

func ParseID(id string) (ID, error) {
	parse, err := uuid.Parse(id)
	if err != nil {
		return ID{}, ErrInvalidID
	}
	return parse, nil
}
