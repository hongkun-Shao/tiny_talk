package db

import (
	"context"
	"time"
	"tiny_talk/utils/config"

	"github.com/go-redis/redis/v8"
)

type RDBClient struct {
	redishandle *redis.Client
}

var RedisClient *RDBClient

func NewRedisClient(ctx *context.Context, rdbCfg *config.RedisConfig) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     rdbCfg.DbAddress,
		Password: rdbCfg.Password, // no password set
		DB:       rdbCfg.Db,       // use default DB
	})

	_, err := rdb.Ping(*ctx).Result()
	if err != nil {
		return err
	}
	RedisClient = &RDBClient{
		redishandle: rdb,
	}
	return nil
}

func Set(ctx *context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return RedisClient.redishandle.Set(*ctx, key, value, expiration)
}

func Del(ctx *context.Context, keys ...string) *redis.IntCmd {
	return RedisClient.redishandle.Del(*ctx, keys...)
}

func Get(ctx *context.Context, key string) *redis.StringCmd {
	return RedisClient.redishandle.Get(*ctx, key)
}

func Expire(ctx *context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return RedisClient.redishandle.Expire(*ctx, key, expiration)
}
