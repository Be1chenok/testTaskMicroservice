package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
	"github.com/golang-jwt/jwt"
)

const (
	passwordSalt = "retd7fg4sd2ud213sgfh"
	signingKey   = "secret"
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

// TODO

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

func (s *AuthentificationService) ParseToken(accesToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

/*
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

	accesToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	if err := s.tokenRepo.SetToken(nil, accesToken, user.Id, tokenTTL); err != nil {
		return "", err
	}

	return accesToken, nil
}

func (s *AuthentificationService) ParseToken(accesToken string) (int, error) {

	userId, err := s.tokenRepo.GetToken(nil, accesToken)
	if err != nil {
		return 0, err
	}

	token, err := jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	if userId != claims.UserId {
		return 0, err
	}

	return claims.UserId, nil
}
*/

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(passwordSalt)))
}
