package domain

type User struct {
	Id       int `db:"id"`
	Email    string
	Username string
	Password string
}
