package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository/postgres"
	rdb "github.com/Be1chenok/testTaskMicroservice/internal/repository/redis"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	PostgresUser
	RedisToken
}

type RedisToken interface {
	SetToken(ctx context.Context, accesToken string, userId int, expiration time.Duration) error
	GetToken(ctx context.Context, accesToken string) (interface{}, error)
}

type PostgresUser interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, passwordHash string) (domain.User, error)
}

func New(db *sql.DB, client *redis.Client) *Repository {
	return &Repository{
		PostgresUser: postgres.NewUser(db),
		RedisToken:   rdb.NewToken(client),
	}
}
