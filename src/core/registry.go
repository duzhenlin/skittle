// Package service_providers
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 10:57

package core

import (
	"github.com/duzhenlin/skittle/v2/src/core/service_providers"
	"github.com/duzhenlin/skittle/v2/src/core/service_providers/interface"
	"github.com/duzhenlin/skittle/v2/src/user/user_service"
	"go.uber.org/dig"
)

var Providers []_interface.ServiceProviderInterface

func RegisterProvider(provider _interface.ServiceProviderInterface) {
	Providers = append(Providers, provider)
}

// RegisterAll 注册所有已注册的Provider（兼容旧版本）
func RegisterAll(container *dig.Container) error {
	return RegisterWithOptions(container, DefaultAppOptions())
}

// RegisterWithOptions 根据选项注册服务
func RegisterWithOptions(container *dig.Container, opts *AppOptions) error {
	// 用户模块（可选）- 必须先注册，因为 log 模块可能依赖它
	if opts.EnableUserModule {
		if err := (&service_providers.UserServiceProvider{}).Register(container); err != nil {
			return err
		}
	} else {
		// 如果 user 模块未启用，但 log 模块需要 user service，提供一个 nil 的 user service
		// 这样 log 模块可以在 user 模块未启用时也能工作
		if opts.EnableLogModule {
			// 提供一个 nil 的 user service，类型必须匹配
			if err := container.Provide(func() *user_service.UserService {
				return nil
			}, dig.Name("user")); err != nil {
				// 如果已经注册了 user，忽略错误（可能在其他地方已注册）
			}
		}
	}

	// 缓存模块
	if opts.EnableCacheModule {
		if err := (&service_providers.CacheServiceProvider{}).Register(container); err != nil {
			return err
		}
	}

	// 日志模块
	if opts.EnableLogModule {
		if err := (&service_providers.LogServiceProvider{}).Register(container); err != nil {
			return err
		}
	}

	// Hprose模块
	if opts.EnableHproseModule {
		if err := (&service_providers.HproseServiceProvider{}).Register(container); err != nil {
			return err
		}
	}

	return nil
}
