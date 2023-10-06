package repository

import (
	"database/sql"

	"github.com/Be1chenok/testTaskMicroservice/internal/repository/postgres"
	rdb "github.com/Be1chenok/testTaskMicroservice/internal/repository/redis"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	PostgresUser postgres.User
	RedisUser    rdb.User
}

func New(db *sql.DB, client *redis.Client) *Repository {
	return &Repository{
		PostgresUser: postgres.NewUser(db),
		RedisUser:    rdb.NewUser(client),
	}
}
