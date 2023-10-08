package model

import "github.com/golang-jwt/jwt"

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpInput struct {
	Email    string
	Username string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	UserId int
}
