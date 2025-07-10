package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/constant"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

var (
	redisClient *redis.Client
	redisOnce   sync.Once
	initErr     error
)

// NewRedisCache 获取RedisCache单例（线程安全）
func NewRedisCache(config *config.Config) (*RedisCache, error) {
	redisOnce.Do(func() {
		redisClient, initErr = initRedis(config)
	})
	if initErr != nil {
		return nil, initErr
	}
	return &RedisCache{client: redisClient}, nil
}

func initRedis(config *config.Config) (*redis.Client, error) {
	if config == nil || config.Redis == nil {
		return nil, errors.New("redis配置不能为空")
	}
	password := config.Redis.Pwd
	db := config.Redis.Db
	poolSize := config.Redis.PoolSize

	switch config.Redis.ConType {
	case constant.RedisConTypeTpc:
		if config.Redis.TpcConArr == "" {
			return nil, errors.New("TCP连接地址不能为空")
		}
		return createTCPClient(config.Redis.TpcConArr, password, db, poolSize)
	case constant.RedisConTypeSentinel:
		if config.Redis.MasterName == "" {
			return nil, errors.New("sentinel模式需要指定MasterName")
		}
		if len(config.Redis.SentinelAddr) == 0 {
			return nil, errors.New("sentinel地址列表不能为空")
		}
		return createSentinelClient(
			config.Redis.MasterName,
			config.Redis.SentinelAddr,
			password,
			db,
		)
	default:
		return nil, fmt.Errorf("不支持的Redis连接类型: %s", config.Redis.ConType)
	}
}

func createTCPClient(addr, password string, db int, poolSize int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})
	return validateConnection(client)
}

func createSentinelClient(masterName string, sentinelAddrs []string, password string, db int) (*redis.Client, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: sentinelAddrs,
		Password:      password,
		DB:            db,
	})
	return validateConnection(client)
}

func validateConnection(client *redis.Client) (*redis.Client, error) {
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis连接验证失败: %w", err)
	}
	return client, nil
}

//  实现 Cache 接口

func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *RedisCache) Set(key string, value string, ttl int) error {
	return r.client.Set(context.Background(), key, value, time.Duration(ttl)*time.Second).Err()
}

func (r *RedisCache) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisCache) Exists(key string) (bool, error) {
	res, err := r.client.Exists(context.Background(), key).Result()
	return res > 0, err
}
func (r *RedisCache) Expire(key string, ttl int) error {
	return r.client.Expire(context.Background(), key, time.Duration(ttl)*time.Second).Err()
}
