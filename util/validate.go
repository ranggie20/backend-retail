package util

import (
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 80

func PasswordValidator(password string) error {
	return passwordvalidator.Validate(password, minEntropyBits)
}
