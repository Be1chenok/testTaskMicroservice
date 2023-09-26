package redis

import "github.com/redis/go-redis/v9"

type User struct {
	cache *redis.Client
}

func NewUser(cache *redis.Client) *User {
	return &User{cache: cache}
}
