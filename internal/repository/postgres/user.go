package postgres

import (
	"database/sql"
	"sync"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
)

type User struct {
	db    *sql.DB
	mutex *sync.Mutex
}

func NewUser(db *sql.DB) *User {
	return &User{
		db:    db,
		mutex: &sync.Mutex{},
	}
}

func (r *User) CreateUser(user domain.User) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var id int

	query := `INSERT INTO users (email, username, password_hash) values ($1, $2, $3) RETURNING id`
	row := r.db.QueryRow(query, user.Email, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *User) GetUserId(username, password string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var userId int

	query := `SELECT id FROM users WHERE username=$1 AND password_hash=$2`
	row := r.db.QueryRow(query, username, password)
	err := row.Scan(&userId)

	return userId, err
}

func (r *User) GetUserIdByAccessToken(accessToken string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var userId int

	query := `SELECT user_id FROM tokens WHERE access_token=$1`
	row := r.db.QueryRow(query, accessToken)
	err := row.Scan(&userId)

	return userId, err
}
func (r *User) GetUserIdByRefreshToken(token string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var userId int

	query := `SELECT user_id FROM tokens WHERE refresh_token=$1`
	row := r.db.QueryRow(query, token)
	err := row.Scan(&userId)

	return userId, err
}

func (r *User) SetTokens(userId int, accessToken, refreshToken string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query := `INSERT INTO tokens (user_id, access_token, refresh_token) values ($1,$2,$3)`

	return r.db.QueryRow(query, userId, accessToken, refreshToken).Err()
}

func (r *User) DeleteUserIdByAccessToken(accessToken string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query := `DELETE FROM tokens WHERE access_token=$1`

	return r.db.QueryRow(query, accessToken).Err()
}
func (r *User) DeleteUserIdByRefreshToken(refreshToken string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query := `DELETE FROM tokens WHERE refresh_token=$1`

	return r.db.QueryRow(query, refreshToken).Err()
}

func (r *User) DeleteAllTokensByUserId(userId int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query := `DELETE FROM tokens WHERE user_id=$1`

	return r.db.QueryRow(query, userId).Err()
}
