package redis

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	client *redis.Client
	mutex  *sync.Mutex
}

func NewUser(client *redis.Client) *User {
	return &User{
		client: client,
		mutex:  &sync.Mutex{},
	}
}

func (r *User) SetAccessToken(ctx context.Context, accessToken string, userId int, expiration time.Duration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.client.Set(ctx, accessToken, userId, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (r *User) GetUserIdByAccessToken(ctx context.Context, accessToken string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	userId, err := r.client.Get(ctx, accessToken).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, errors.New("invalid or expired token")
		}
		return 0, err
	}

	return userId, nil
}

func (r *User) DeleteUserIdByAccessToken(ctx context.Context, accessToken string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.client.Del(ctx, accessToken).Err(); err != nil {
		return errors.New("invalid or expire token")
	}

	return nil
}

func (r *User) DeleteAllAccessTokensByUserId(ctx context.Context, userId int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	keys, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		return err
	}
	var deletedKeys int
	for _, key := range keys {
		value, err := r.client.Get(ctx, key).Int()
		if err != nil {
			return err
		}

		if value == userId {
			if err := r.client.Del(ctx, key).Err(); err != nil {
				return err
			}
			deletedKeys++
		}
	}

	return nil

}
