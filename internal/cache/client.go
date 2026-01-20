package cache

import (
	"cluster-agent/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		PoolSize:     100,
		MinIdleConns: 10,

		PoolTimeout: 4 * time.Second,

		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	cleanup := func() {
		log.Println("Closing Redis connection")
		if err := rdb.Close(); err != nil {
			log.Printf("Error closing redis: %v", err)
		}
	}

	return rdb, cleanup, nil
}
