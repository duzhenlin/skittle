// Package core
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:17

package core

import (
	"context"

	"github.com/duzhenlin/skittle/v2/src/cache"
	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/hprose/hprose_client"
	"github.com/duzhenlin/skittle/v2/src/hprose/hprose_server"
	"github.com/duzhenlin/skittle/v2/src/log"
	"github.com/duzhenlin/skittle/v2/src/user/user_auth"
	"github.com/duzhenlin/skittle/v2/src/user/user_cache"
	"github.com/duzhenlin/skittle/v2/src/user/user_client"
	"github.com/duzhenlin/skittle/v2/src/user/user_service"
	"go.uber.org/dig"
)

// App 应用核心结构
// User相关字段使用指针，支持可选模块
// 使用 NewAppWithDefaults() 可快速创建包含所有模块的应用（即开即用）
// 使用 NewCoreAppOnly() 可创建仅核心功能的应用（不包含user模块）
type App struct {
	Version    string
	Config     *config.Config
	Ctx        context.Context
	Client     *hprose_client.HproseClientService
	Server     *hprose_server.HproseServerService
	Auth       *user_auth.AuthService      // 可选：用户认证服务
	User       *user_service.UserService   // 可选：用户服务
	UserClient *user_client.UserClientService // 可选：用户客户端
	UserCache  *user_cache.UserCacheService   // 可选：用户缓存
	Cache      cache.Cache
	Logger     log.Logger
}

// App基础依赖模块（完整模式）
// 注意：dig 不支持 optional 标签，我们通过条件注册来控制依赖
type appDeps struct {
	dig.In
	Cfg          *config.Config
	Ctx          context.Context
	HproseClient *hprose_client.HproseClientService `name:"hprose_client"`
	HproseServer *hprose_server.HproseServerService `name:"hprose_server"`
	Auth         *user_auth.AuthService             `name:"auth"`
	User         *user_service.UserService          `name:"user"`
	UserCache    *user_cache.UserCacheService       `name:"user_cache"`
	UserClient   *user_client.UserClientService      `name:"user_client"`
	Cache        cache.Cache                        `name:"cache"`
	Logger       log.Logger                         `name:"log"`
}

// App核心依赖模块（不包含user）
// 注意：dig 不支持 optional 标签，我们通过条件注册来控制依赖
type coreAppDeps struct {
	dig.In
	Cfg          *config.Config
	Ctx          context.Context
	HproseClient *hprose_client.HproseClientService `name:"hprose_client"`
	HproseServer *hprose_server.HproseServerService `name:"hprose_server"`
	Cache        cache.Cache                        `name:"cache"`
	Logger       log.Logger                         `name:"log"`
}

// NewApp 创建应用实例（完整模式，包含user模块）
// 这是默认的即开即用方式
func NewApp(c *dig.Container) (*App, error) {
	var app *App
	err := c.Invoke(func(deps appDeps) {
		app = &App{
			Version:    GetVersion(),
			Config:     deps.Cfg,
			Ctx:        deps.Ctx,
			Client:     deps.HproseClient,
			Server:     deps.HproseServer,
			User:       deps.User,
			Auth:       deps.Auth,
			UserClient: deps.UserClient,
			UserCache:  deps.UserCache,
			Cache:      deps.Cache,
			Logger:     deps.Logger,
		}
	})
	if err != nil {
		return nil, err
	}
	return app, nil
}

// NewCoreApp 创建核心应用实例（不包含user模块）
// 适用于只需要核心功能的场景
func NewCoreApp(c *dig.Container) (*App, error) {
	var app *App
	err := c.Invoke(func(deps coreAppDeps) {
		app = &App{
			Version:    GetVersion(),
			Config:     deps.Cfg,
			Ctx:        deps.Ctx,
			Client:     deps.HproseClient,
			Server:     deps.HproseServer,
			Cache:      deps.Cache,
			Logger:     deps.Logger,
			// User相关字段保持为nil
		}
	})
	if err != nil {
		return nil, err
	}
	return app, nil
}

// HasUserModule 检查是否启用了用户模块
func (a *App) HasUserModule() bool {
	return a.User != nil && a.Auth != nil
}
