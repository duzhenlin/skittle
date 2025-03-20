// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 18:37

package user

import (
	"context"
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/hprose/client"
	"github.com/go-redis/redis/v8"
	"time"
)

// User 用户服务结构体
type User struct {
	ctx          context.Context
	config       *config.Config
	redisClient  *redis.Client  // 注入的Redis客户端
	hproseClient *client.Client // Hprose客户端
}

// NewUser 创建用户服务实例
func NewUser(ctx context.Context, cfg *config.Config) *User {
	return &User{
		ctx:    ctx,
		config: cfg,
	}
}

// SetHproseClient 设置hprose客户端
func (u *User) SetHproseClient(client *client.Client) *User {
	u.hproseClient = client
	return u
}

// SetRedisClient 设置redis客户端
func (u *User) SetRedisClient(redisClient *redis.Client) *User {
	u.redisClient = redisClient
	return u
}

func (u *User) GetConfig() *config.Config {
	return u.config
}

// CacheKey 生成缓存键
func (u *User) CacheKey(suffix string) string {
	return fmt.Sprintf("%s:%s", u.config.Skittle.Namespace, suffix)
}

func (u *User) SetUserInfo(userInfo *LoginData) (string, error) {
	// 设置缓存
	key := u.CacheKey("u_" + userInfo.ModuleToken)
	result, err := u.redisClient.Set(u.ctx, key, userInfo, 2*time.Hour).Result()
	if err != nil {
		return "", fmt.Errorf("用户缓存写入失败: %w", err)
	}
	return result, err
}

func (u *User) GetUserInfo(token string) (interface{}, error) {
	key := u.CacheKey("u_" + token)
	result, err := u.redisClient.Get(u.ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("用户缓存读取失败: %w", err)
	}
	return result, err
}
