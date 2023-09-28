package service

import (
	"context"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
)

type Service struct {
	Authentification
}

type Authentification interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, token string) (int, error)
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Authentification: NewAuthentificationService(repo.PostgresUser, repo.RedisToken),
	}
}
