package repository

import (
	"database/sql"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository/postgres"
	rdb "github.com/Be1chenok/testTaskMicroservice/internal/repository/redis"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Authentification
	Validation
}

type Validation interface {
}

type Authentification interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

func New(db *sql.DB, cache *redis.Client) *Repository {
	return &Repository{
		Authentification: postgres.NewUser(db),
		Validation:       rdb.NewUser(cache),
	}
}
