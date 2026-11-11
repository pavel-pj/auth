package tests

import (
	"context"
	"testing"

	"auth/internal/repository"

	"github.com/gojuno/minimock/v3"
)

func TestAuthLogin(t *testing.T) {

	type AuthRepositoryTestFunc func(mc *minimock.Controller) repository.AuthRepository

	type args struct {
		ctx      context.Context
		login    string
		password string
	}

	var (
	//login    = "oleg"
	//password = "qwerty206"
	)

}
