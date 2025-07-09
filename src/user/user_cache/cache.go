// Package user_cache
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/7
// @Time: 19:44

// Package user_cache 用户缓存相关
package user_cache

import (
	"context"

	"github.com/duzhenlin/skittle/v2/src/cache"
	"github.com/duzhenlin/skittle/v2/src/config"
	"go.uber.org/dig"
)

type UserCacheService struct {
	Ctx    context.Context
	cache  cache.Cache
	config *config.Config
}

type CacheDeps struct {
	dig.In
	Ctx    context.Context
	Config *config.Config
	Cache  cache.Cache `name:"cache"`
}

func NewUserCache(deps CacheDeps) *UserCacheService {
	return &UserCacheService{
		Ctx:    deps.Ctx,
		config: deps.Config,
		cache:  deps.Cache,
	}
}
