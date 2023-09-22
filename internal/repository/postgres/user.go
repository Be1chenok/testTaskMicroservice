package postgres

import (
	"database/sql"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

func (r *User) CreateUser(user domain.User) (uint, error) {
	var id uint
	return id, nil
}
