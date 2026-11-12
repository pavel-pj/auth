package utils

import (
	"github.com/sppps/pwdhash"
)

func VerifyPassword(password, hash string) bool {
	return pwdhash.Validate(password, hash) == nil
}
