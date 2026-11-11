package model

import "github.com/dgrijalva/jwt-go"

type UserLogin struct {
	Name     string
	Password string
	Role     string
}

type UserAuthInfo struct {
	Name string `json:"username"`
	Role string `json:"role"`
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
