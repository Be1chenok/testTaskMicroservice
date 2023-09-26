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

/*
func (rdb *RedisDB) Set(ctx context.Context, key string, value interface{}) bool {
	return true
}

func (rdb *RedisDB) Get(ctx context.Context, key string) interface{} {
	return nil
}

func (rdb *RedisDB) Delete(ctx context.Context, key string) interface{} {
	return rdb.client.Del(ctx, key)
}
*/
