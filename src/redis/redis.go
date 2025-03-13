// Package redis
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 22:27

package redis

import (
	"github.com/duzhenlin/skittle/src/config"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	config *config.Config
}

var client *redis.Client

func GetClientInstance(config *config.Config) *redis.Client {
	if client == nil {
		client = InitRedis(*config)
	}
	return client
}

// InitRedis 初始化redis连接
func InitRedis(config config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Redis.Cont,
		Password: config.Redis.Pwd, // no password set
		DB:       config.Redis.Db,  // use default DB
	})
}
