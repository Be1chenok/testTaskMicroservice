package model

import "github.com/golang-jwt/jwt"

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Tokens struct {
	AccesToken   string
	RefreshToken string
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	UserId int
}
