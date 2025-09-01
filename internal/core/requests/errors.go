package requests

import "errors"

var (
	ErrPhoneDigitsCount  = errors.New("неправильное количество цифр в телефонном номере")
	ErrPhoneWrongSymbols = errors.New("телефонный номер содержит неподдерживаемые символы")
)
