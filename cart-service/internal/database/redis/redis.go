package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sabirkekw/ecommerce_go/cart-service/internal/config"
)

func ConnectToRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%v", cfg.Redis.Host, cfg.Redis.Port),
		DB:   0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("failed to ping redis database")
	}
	return client
}
