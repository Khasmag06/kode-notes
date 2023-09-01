package redis

import (
	"context"
	"fmt"
	"github.com/Khasmag06/kode-notes/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context, redisConfig config.RedisConfig) (*redis.Client, error) {
	redisAddress := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf(`connect redis ping: %w`, err)
	}

	return rdb, nil
}
