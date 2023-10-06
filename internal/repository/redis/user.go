package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type User interface {
	SetAccessToken(ctx context.Context, accesToken string, userId int, expiration time.Duration) error

	GetUserIdByAccessToken(ctx context.Context, accesToken string) (int, error)

	DeleteUserIdByAccessToken(ctx context.Context, accesToken string) error
	DeleteAllAccessTokensByUserId(ctx context.Context, userId int) (int, error)
}

type UserRepo struct {
	client *redis.Client
}

func NewUser(client *redis.Client) *UserRepo {
	return &UserRepo{
		client: client,
	}
}

func (r *UserRepo) SetAccessToken(ctx context.Context, accessToken string, userId int, expiration time.Duration) error {
	if err := r.client.Set(ctx, accessToken, userId, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUserIdByAccessToken(ctx context.Context, accessToken string) (int, error) {
	userId, err := r.client.Get(ctx, accessToken).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, fmt.Errorf("Invalid or expired token")
		}
		return 0, err
	}

	return userId, nil
}

func (r *UserRepo) DeleteUserIdByAccessToken(ctx context.Context, accessToken string) error {
	if err := r.client.Del(ctx, accessToken).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("Invalid or expired token")
		}
		return err
	}

	return nil
}

func (r *UserRepo) DeleteAllAccessTokensByUserId(ctx context.Context, userId int) (int, error) {
	keys, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		return 0, err
	}
	var deletedKeys int
	for _, key := range keys {
		value, err := r.client.Get(ctx, key).Int()
		if err != nil {
			return deletedKeys, err
		}

		if value == userId {
			if err := r.client.Del(ctx, key).Err(); err != nil {
				if errors.Is(err, redis.Nil) {
					return deletedKeys, fmt.Errorf("Invalid or expired token")
				}
				return deletedKeys, err
			}
			deletedKeys++
		}
	}

	return deletedKeys, nil
}
