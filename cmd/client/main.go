package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:8848",
		DB:   0,
	})
	err := client.Get(context.Background(), "TEST").Err()
	if err != nil {
		panic(err)
	}
}
