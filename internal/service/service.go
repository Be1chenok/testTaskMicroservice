package service

import (
	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
)

type Service struct {
	Authentification
	Validation
}

type Authentification interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
}

type Validation interface {
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Authentification: NewAuthentificationService(repo.Authentification),
		Validation:       NewValidationService(repo.Validation),
	}
}
