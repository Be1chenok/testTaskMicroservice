package service

import (
	"crypto/sha256"
	"fmt"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
)

const passwordSalt = "retd7fg4sd2ud213sgfh"

type AuthentificationService struct {
	repo repository.Authentification
}

func NewAuthentificationService(repo repository.Authentification) *AuthentificationService {
	return &AuthentificationService{repo: repo}
}

func (s *AuthentificationService) CreateUser(user domain.User) (uint, error) {
	user.Password = genPasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func genPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(passwordSalt)))
}
