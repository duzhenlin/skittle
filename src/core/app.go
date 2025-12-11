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
	"github.com/duzhenlin/skittle/v2/src/core/service_providers"
	"github.com/duzhenlin/skittle/v2/src/hprose/hprose_client"
	"github.com/duzhenlin/skittle/v2/src/hprose/hprose_server"
	"github.com/duzhenlin/skittle/v2/src/log"
	"github.com/duzhenlin/skittle/v2/src/user/user_auth"
	"github.com/duzhenlin/skittle/v2/src/user/user_cache"
	"github.com/duzhenlin/skittle/v2/src/user/user_client"
	"github.com/duzhenlin/skittle/v2/src/user/user_service"
	"go.uber.org/dig"
)

type App struct {
	Version    string
	Config     *config.Config
	Ctx        context.Context
	Client     *hprose_client.HproseClientService
	Server     *hprose_server.HproseServerService
	Auth       *user_auth.AuthService
	User       *user_service.UserService
	UserClient *user_client.UserClientService
	UserCache  *user_cache.UserCacheService
	Cache      cache.Cache
	Logger     log.Logger
}

// App基础依赖模块
type appDeps struct {
	dig.In
	Cfg          *config.Config
	Ctx          context.Context
	HproseClient *hprose_client.HproseClientService `name:"hprose_client"`
	HproseServer *hprose_server.HproseServerService `name:"hprose_server"`
	Auth         *user_auth.AuthService             `name:"auth"`
	User         *user_service.UserService          `name:"user"`
	UserCache    *user_cache.UserCacheService       `name:"user_cache"`
	UserClient   *user_client.UserClientService     `name:"user_client"`
	Cache        cache.Cache                        `name:"cache"`
	Logger       log.Logger                         `name:"log"`
}

func init() {
	// 注册基础依赖模块
	RegisterProvider(&service_providers.CacheServiceProvider{})
	RegisterProvider(&service_providers.UserServiceProvider{})
	RegisterProvider(&service_providers.HproseServiceProvider{})
	RegisterProvider(&service_providers.LogServiceProvider{})
}

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
