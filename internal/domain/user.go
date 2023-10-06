package domain

type User struct {
	Id       int    `json:"id,omitempty"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
