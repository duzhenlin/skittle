// Package redis
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 22:27

package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/constant"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
	initErr     error
)

// GetRedisClient 获取Redis客户端单例（线程安全）
func GetRedisClient(config *config.Config) (*redis.Client, error) {
	redisOnce.Do(func() {
		redisClient, initErr = initRedis(config)
	})
	return redisClient, initErr
}

// initRedis 封装Redis初始化逻辑
func initRedis(config *config.Config) (*redis.Client, error) {
	// 参数校验
	if config == nil || config.Redis == nil {
		return nil, errors.New("redis配置不能为空")
	}

	// 公共配置项
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

// createTCPClient 创建TCP直连客户端
func createTCPClient(addr, password string, db int, poolSize int) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: poolSize, // 设置最大连接数
	})
	return validateConnection(client)
}

// createSentinelClient 创建Sentinel模式客户端
func createSentinelClient(masterName string, sentinelAddrs []string, password string, db int) (*redis.Client, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: sentinelAddrs,
		Password:      password,
		DB:            db,
	})
	return validateConnection(client)
}

// validateConnection 验证Redis连接有效性
func validateConnection(client *redis.Client) (*redis.Client, error) {
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis连接验证失败: %w", err)
	}
	return client, nil
}
