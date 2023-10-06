package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Be1chenok/testTaskMicroservice/internal/domain"
)

type User interface {
	CreateUser(user domain.User) (int, error)

	SetTokens(userId int, accessToken, refreshToken string) error

	GetUserId(username, passwordHash string) (int, error)
	GetUserIdByAccessToken(accessToken string) (int, error)
	GetUserIdByRefreshToken(refreshToken string) (int, error)

	DeleteUserIdByAccessToken(accessToken string) error
	DeleteUserIdByRefreshToken(refreshToken string) error
	DeleteAllTokensByUserId(userId int) error
}

type UserRepo struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) CreateUser(user domain.User) (int, error) {
	var id int

	query := `INSERT INTO users (email, username, password_hash) values ($1, $2, $3) RETURNING id`
	row := r.db.QueryRow(query, user.Email, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepo) GetUserId(username, password string) (int, error) {
	var userId int

	row := r.db.QueryRow(
		`SELECT id FROM users WHERE username=$1 AND password_hash=$2`,
		username, password,
	)
	if err := row.Scan(&userId); err != nil {
		return 0, fmt.Errorf("Failed to find user: %v", err)
	}

	return userId, nil
}

func (r *UserRepo) GetUserIdByAccessToken(accessToken string) (int, error) {
	var userId int

	row := r.db.QueryRow(
		`SELECT user_id FROM tokens WHERE access_token=$1`,
		accessToken,
	)
	if err := row.Scan(&userId); err != nil {
		return 0, fmt.Errorf("Failed to find by access token: %v", err)
	}

	return userId, nil
}
func (r *UserRepo) GetUserIdByRefreshToken(token string) (int, error) {
	var userId int

	row := r.db.QueryRow(
		`SELECT user_id FROM tokens WHERE refresh_token=$1`,
		token,
	)
	if err := row.Scan(&userId); err != nil {
		return 0, fmt.Errorf("Failed to find by refresh token: %v", err)
	}

	return userId, nil
}

func (r *UserRepo) SetTokens(userId int, accessToken, refreshToken string) error {
	_, err := r.db.Exec(
		`INSERT INTO tokens (user_id, access_token, refresh_token) values ($1,$2,$3)`,
		userId,
		accessToken,
		refreshToken,
	)
	if err != nil {
		return fmt.Errorf("Failed to set tokens: %v", err)
	}

	return nil
}

func (r *UserRepo) DeleteUserIdByAccessToken(accessToken string) error {
	_, err := r.db.Exec(
		`DELETE FROM tokens WHERE access_token=$1`,
		accessToken,
	)
	if err != nil {
		return fmt.Errorf("Failed to delete by access token: %v", err)
	}

	return nil
}

func (r *UserRepo) DeleteUserIdByRefreshToken(refreshToken string) error {
	_, err := r.db.Exec(
		`DELETE FROM tokens WHERE refresh_token=$1`,
		refreshToken,
	)
	if err != nil {
		return fmt.Errorf("Failed to delete by refresh token: %v", err)
	}

	return nil
}

func (r *UserRepo) DeleteAllTokensByUserId(userId int) error {
	_, err := r.db.Exec(
		`DELETE FROM tokens WHERE user_id=$1`,
		userId,
	)
	if err != nil {
		return fmt.Errorf("Failed to delete by user id: %v", err)
	}

	return nil
}
