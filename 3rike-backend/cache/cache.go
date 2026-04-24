package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func New(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	c := redis.NewClient(opts)
	return c, c.Ping(context.Background()).Err()
}
