package config

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
	ctx   = context.Background()
)

func InitRedis() error {

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		db = 0
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	_, err = Redis.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	// err = Redis.Set(ctx, "test", "Hello Redis", 0).Err()
	// if err != nil {
	// 	return err
	// }

	// value, err := Redis.Get(ctx, "test").Result()
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(value)

	fmt.Println("✅ Redis Connected Successfully")

	return nil
}
