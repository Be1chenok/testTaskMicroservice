package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
	"github.com/golang-jwt/jwt"
)

const (
	passwordSalt = "retd7fg4sd2ud213sgfh"
	signingKey   = "IPofjasld#ASDsaqQWwe#thnmE#dfcv"
	tokenTTL     = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type AuthentificationService struct {
	userRepo  repository.PostgresUser
	tokenRepo repository.RedisToken
}

func NewAuthentificationService(userRepo repository.PostgresUser, tokenRepo repository.RedisToken) *AuthentificationService {
	return &AuthentificationService{userRepo: userRepo, tokenRepo: tokenRepo}
}

func (s *AuthentificationService) CreateUser(user domain.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.userRepo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(passwordSalt)))
}

func (s *AuthentificationService) GenerateToken(username, password string) (string, error) {
	user, err := s.userRepo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}
