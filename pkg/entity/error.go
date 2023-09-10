package entity

import "errors"

var ErrInvalidEntity = errors.New("invalid entity")
var ErrIDIsRequire = errors.New("id is required")
var ErrInvalidID = errors.New("invalid id")

func NewError(text string) error {
	return errors.New(text)
}
