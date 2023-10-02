package service

import (
	"context"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
	"github.com/Be1chenok/testTaskMicroservice/pkg/hash"
	"github.com/golang-jwt/jwt"
)

type SignInInput struct {
	Username string
	Password string
}

type SignUpInput struct {
	Email    string
	Username string
	Password string
}

type Tokens struct {
	AccesToken   string
	RefreshToken string
}

type accesTokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type Service struct {
	Authentification
}

type Authentification interface {
	SignUp(input SignUpInput) (int, error)
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)

	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	ParseToken(ctx context.Context, accesToken string) (int, error)

	SignOut(ctx context.Context, accessToken string) error
	FullSignOut(ctx context.Context, accesToken string) error
}

func New(repo *repository.Repository, hasher *hash.SHA256Hasher, tokensSigningKey string, accesTokenTTL, refreshTokenTTL time.Duration) *Service {
	return &Service{
		Authentification: NewAuthentificationService(repo.PostgresUser, repo.RedisUser, hasher, tokensSigningKey, accesTokenTTL, refreshTokenTTL),
	}
}
