package handlers

import (
	"context"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func NewCache() (*cache.Cache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	myCache := cache.New(&cache.Options{
		Redis: client,
	})
	return myCache, nil
}
