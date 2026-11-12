package utils

import (
	"auth/internal/model"
	"time"

	"errors"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(user model.UserAuthInfo, sercretKey []byte, duration time.Duration) (string, error) {

	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: user.Name,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(sercretKey)

}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {

			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}

			return secretKey, nil
		},
	)

	if err != nil {
		return nil, err

	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.New("Uncorrect token format")
	}

	return claims, nil

}
