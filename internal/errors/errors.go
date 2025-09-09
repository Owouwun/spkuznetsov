package core_errors

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyField         = errors.New("поле должно быть заполнено")
	ErrInvalidDate        = errors.New("недопустимая дата")
	ErrInvalidPhoneNumber = errors.New("недопустимый номер телефона")

	ErrRequestActionNotPermittedByStatus = errors.New("недопустимое действие для текущего статуса заявки")

	ErrNotImplemented = errors.New("not implemented")
)

func NewErrEmptyField(field string) error {
	return fmt.Errorf("%w: %s", ErrEmptyField, field)
}
