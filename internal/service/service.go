package service

import (
	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
)

type Service struct {
	Registration
}

type Registration interface {
	CreateUser(user domain.User) (uint, error)
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Registration: NewRegistrationService(repo.Registration),
	}
}
