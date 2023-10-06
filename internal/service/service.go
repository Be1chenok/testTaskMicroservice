package service

import (
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/repository"
	"github.com/Be1chenok/testTaskMicroservice/pkg/hash"
)

type Service struct {
	Authentification
}

func New(repo *repository.Repository, hasher *hash.SHA256Hasher, tokensSigningKey string, accesTokenTTL, refreshTokenTTL time.Duration) *Service {
	return &Service{
		Authentification: NewAuthentificationService(repo.PostgresUser, repo.RedisUser, hasher, tokensSigningKey, accesTokenTTL, refreshTokenTTL),
	}
}
