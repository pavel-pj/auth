package auth

import "context"

func (s AuthService) Login(ctx context.Context, login, password string) (string, error) {

	token, err := s.Repository.Login(ctx, login, password)
	if err != nil {
		return "", err
	}

	return token, nil
}
