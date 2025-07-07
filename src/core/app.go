// Package core
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:17

package core

import (
	"context"
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/hprose/client"
	"github.com/duzhenlin/skittle/src/hprose/server"
	"github.com/duzhenlin/skittle/src/log_client"
	redis2 "github.com/duzhenlin/skittle/src/redis"
	"github.com/duzhenlin/skittle/src/user"
	"github.com/go-redis/redis/v8"
)

const Version = "1.1.4"

const (
	ProviderClient = "Client"
	ProviderServer = "Server"
	ProviderUser   = "User"
	ProviderLog    = "Log"
)

// 使用切片明确维护需要初始化的组件列表
var providerList = []string{
	ProviderClient,
	ProviderServer,
	ProviderUser,
	ProviderLog,
}

type App struct {
	Version     string
	Config      *config.Config        // 配置
	Client      *client.Client        // 客户端
	Server      *server.Server        // 服务端
	User        *user.User            // 用户
	Ctx         context.Context       // 上下文
	RedisClient *redis.Client         // redis客户端
	Log         *log_client.LogClient // 日志客户端
}

// NewApp 创建应用实例，返回实例和可能的错误
func NewApp(ctx context.Context, config *config.Config) (*App, error) {
	a := &App{
		Config:  config,
		Ctx:     ctx,
		Version: Version,
	}

	if err := a.RegisterBaseProviders(); err != nil {
		return nil, fmt.Errorf("应用初始化失败: %w", err)
	}

	if err := a.RegisterProviders(); err != nil {
		return nil, fmt.Errorf("应用初始化失败: %w", err)
	}

	return a, nil
}

// 定义组件初始化函数类型
type componentInitializer func(a *App) interface{}

// 组件初始化映射表
var componentInitializers = map[string]componentInitializer{
	ProviderClient: func(a *App) interface{} {
		return client.NewClient(a.Ctx, a.Config)
	},
	ProviderServer: func(a *App) interface{} {
		return server.NewServer(a.Ctx, a.Config)
	},
	ProviderUser: func(a *App) interface{} {
		return user.NewUser(a.Ctx, a.Config).
			SetRedisClient(a.RedisClient).
			SetHproseClient(a.Client)
	},
	ProviderLog: func(a *App) interface{} {
		return log_client.NewLog(a.Ctx, a.Config, a.RedisClient)
	},
}

// RegisterProviders 注册所有依赖组件
func (a *App) RegisterProviders() error {
	for _, provider := range providerList {
		initializer, exists := componentInitializers[provider]
		if !exists {
			return fmt.Errorf("未定义的组件: %s", provider)
		}

		instance := initializer(a)

		// 使用类型断言赋值到对应字段
		switch provider {
		case ProviderClient:
			a.Client = instance.(*client.Client)
		case ProviderServer:
			a.Server = instance.(*server.Server)
		case ProviderUser:
			a.User = instance.(*user.User)
		case ProviderLog:
			a.Log = instance.(*log_client.LogClient)

		}
	}
	return nil
}

func (a *App) RegisterBaseProviders() error {
	redisClient, err := redis2.GetRedisClient(a.Config)
	if err != nil {
		return fmt.Errorf("获取redis客户端失败: %w", err)
	}
	a.RedisClient = redisClient
	return nil
}
