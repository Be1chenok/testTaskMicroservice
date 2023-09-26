package redis

import (
	"context"

	"github.com/Be1chenok/testTaskMicroservice/internal/config"
	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, conf *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Cache.Host + ":" + conf.Cache.Port,
		Password: conf.Cache.Password,
		DB:       conf.Cache.DB,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
