package auth

import "context"

func (r *AuthRepository) Login(ctx context.Context, login, password string) (string, error) {

	return "superPuperToken", nil

}
