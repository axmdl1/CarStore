package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func NewClient(addr, pass string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
}
