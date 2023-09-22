package service

import (
	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
)

type RegistrationService struct {
	repo repository.Registration
}

func NewRegistrationService(repo repository.Registration) *RegistrationService {
	return &RegistrationService{repo: repo}
}

func (s *RegistrationService) CreateUser(user domain.User) (uint, error) {
	return s.repo.CreateUser(user)
}
