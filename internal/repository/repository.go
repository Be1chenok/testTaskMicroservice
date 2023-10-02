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
	RedisUser
}

type RedisUser interface {
	SetAccessToken(ctx context.Context, accesToken string, userId int, expiration time.Duration) error
	GetUserIdByAccessToken(ctx context.Context, accesToken string) (int, error)
	DeleteUserIdByAccessToken(ctx context.Context, accesToken string) error
	DeleteAllAccessTokensByUserId(ctx context.Context, userId int) error
}

type PostgresUser interface {
	CreateUser(user domain.User) (int, error)

	SetTokens(userId int, accessToken, refreshToken string) error

	GetUserId(username, passwordHash string) (int, error)
	GetUserIdByAccessToken(accessToken string) (int, error)
	GetUserIdByRefreshToken(refreshToken string) (int, error)

	DeleteUserIdByAccessToken(accessToken string) error
	DeleteUserIdByRefreshToken(refreshToken string) error
	DeleteAllTokensByUserId(userId int) error
}

func New(db *sql.DB, client *redis.Client) *Repository {
	return &Repository{
		PostgresUser: postgres.NewUser(db),
		RedisUser:    rdb.NewUser(client),
	}
}
