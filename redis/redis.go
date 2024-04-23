package cache

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var Redis *redis.Client

func Init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})
}

func Get(c *gin.Context, key string) (string, error) {
	val, redisErr := Redis.Get(c.Request.Context(), key).Result()
	if redisErr == redis.Nil {
		return "", errors.New(redisErr.Error())
	}
	return val, nil
}

func Set(c *gin.Context, key, value string, duration time.Duration) error {
	redisErr := Redis.Set(c.Request.Context(), key, value, duration).Err()
	if redisErr != nil {
		return errors.New(redisErr.Error())
	}
	return nil
}
