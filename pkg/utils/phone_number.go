package utils

import (
	"errors"
	"regexp"
	"strings"
)

func isContainsNumbersOnly(phoneNumber string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(phoneNumber)
}

func getPhoneWithoutDecorators(phoneNumber string) string {
	phoneWithoutSpaces := strings.ReplaceAll(phoneNumber, " ", "")
	phoneWithoutLeftParentheses := strings.ReplaceAll(phoneWithoutSpaces, "(", "")
	phoneWithoutRightParentheses := strings.ReplaceAll(phoneWithoutLeftParentheses, ")", "")
	phoneWithoutHyphen := strings.ReplaceAll(phoneWithoutRightParentheses, "-", "")
	phoneWithoutPlus := strings.ReplaceAll(phoneWithoutHyphen, "+", "")

	return phoneWithoutPlus
}

func StandartizePhoneNumber(phoneNumber string) (string, error) {
	clearPhone := getPhoneWithoutDecorators(phoneNumber)

	if !isContainsNumbersOnly(clearPhone) {
		return "", errors.New("номер телефона содержит посторонние символы")
	}

	phoneLen := len(clearPhone)
	if phoneLen == 7 {
		return clearPhone, nil
	}
	if phoneLen == 11 {
		if clearPhone[0] == '8' {
			clearPhone = "7" + clearPhone[1:]
		}
		return "+" + clearPhone, nil
	}

	return "", errors.New("номер телефона имеет некорректное количество цифр")
}
