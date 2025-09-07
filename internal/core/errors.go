package core

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyField  = errors.New("поле должно быть заполнено")
	ErrInvalidDate = errors.New("недопустимая дата")

	ErrRequestActionNotPermittedByStatus = errors.New("недопустимое действие для текущего статуса заявки")

	ErrNotImplemented = errors.New("not implemented")
)

func NewErrEmptyField(field string) error {
	return fmt.Errorf("%w: %s", ErrEmptyField, field)
}
