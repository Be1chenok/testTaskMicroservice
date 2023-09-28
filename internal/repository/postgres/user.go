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

func (r *User) CreateUser(user domain.User) (int, error) {
	var id int
	query := `INSERT INTO users (email, username, password_hash) values ($1, $2, $3) RETURNING id`
	row := r.db.QueryRow(query, user.Email, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *User) GetUser(username, password string) (domain.User, error) {
	var user domain.User

	query := `SELECT id FROM users WHERE username=$1 AND password_hash=$2`
	row := r.db.QueryRow(query, username, password)
	err := row.Scan(&user.Id)

	return user, err

}
