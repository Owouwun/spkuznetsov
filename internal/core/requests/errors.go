package requests

import "errors"

var (
	ErrPhoneDigitsCount           = errors.New("неправильное количество цифр в телефонном номере")
	ErrPhoneWrongSymbols          = errors.New("телефонный номер содержит неподдерживаемые символы")
	ErrActionNotPermittedByStatus = errors.New("недопустимое действие для текущего статуса заявки")

	ErrNotImplemented = errors.New("not implemented")
)
