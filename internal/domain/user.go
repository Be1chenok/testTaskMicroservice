package domain

import (
	"time"
)

type User struct {
	Id           int
	Email        string
	Username     string
	Password     string
	RegisteresAt time.Time
}
