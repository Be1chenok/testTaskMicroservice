package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Token struct {
	client *redis.Client
}

func NewToken(client *redis.Client) *Token {
	return &Token{client: client}
}

func (r *Token) SetToken(ctx context.Context, accesToken string, userId int, expiration time.Duration) error {
	if err := r.client.Set(ctx, accesToken, userId, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Token) GetToken(ctx context.Context, accesToken string) (*int, error) {
	userId, err := r.client.Get(ctx, accesToken).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("invalid or expired token")
		}
		return nil, err
	}
	return &userId, nil
}
