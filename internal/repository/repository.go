package repository

import (
	"database/sql"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
	"github.com/Be1chenok/testTaskMicroservice/internal/repository/postgres"
)

type Repository struct {
	Authentification
}

type Authentification interface {
	CreateUser(user domain.User) (uint, error)
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Authentification: postgres.NewUser(db),
	}
}
