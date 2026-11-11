package utils

import (
	"github.com/sppps/pwdhash"
)

func VerifyPassword(hashedPassword, candidatePassword string) bool {

	err := pwdhash.Validate(hashedPassword, candidatePassword)
	return err == nil
}
