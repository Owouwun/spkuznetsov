package utils

import (
	"errors"
	"time"
)

func MustNotPast(date *time.Time) error {
	if date.Before(time.Now()) {
		return errors.New("can't be past")
	}
	return nil
}
