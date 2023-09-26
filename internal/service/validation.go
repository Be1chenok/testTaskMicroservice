package service

import "github.com/Be1chenok/testTaskMicroservice/internal/repository"

type ValidationService struct {
	repo repository.Validation
}

func NewValidationService(repo repository.Validation) *ValidationService {
	return &ValidationService{repo: repo}
}
