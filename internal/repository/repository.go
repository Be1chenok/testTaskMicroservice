package repository

import (
	"database/sql"

	"github.com/Be1chenok/testTaskMicroservice/internal/repository/postgres"
	rdb "github.com/Be1chenok/testTaskMicroservice/internal/repository/redis"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Postgres postgres.PGUser
	Redis    rdb.RUser
}

func New(db *sql.DB, client *redis.Client) *Repository {
	return &Repository{
		Postgres: postgres.NewUser(db),
		Redis:    rdb.NewUser(client),
	}
}
