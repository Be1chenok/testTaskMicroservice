package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Token struct {
	client *redis.Client
}

func NewToken(client *redis.Client) *Token {
	return &Token{client: client}
}

func (r *Token) Set(ctx context.Context, key string, value interface{}) bool {
	return true
}

func (r *Token) Get(ctx context.Context, key string) interface{} {
	return nil
}
