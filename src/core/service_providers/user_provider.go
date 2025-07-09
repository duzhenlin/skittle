// Package service_providers
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 10:55

package service_providers

import (
	"github.com/duzhenlin/skittle/v2/src/core/service_providers/interface"
	"github.com/duzhenlin/skittle/v2/src/user/user_auth"
	"github.com/duzhenlin/skittle/v2/src/user/user_cache"
	"github.com/duzhenlin/skittle/v2/src/user/user_client"
	"github.com/duzhenlin/skittle/v2/src/user/user_service"
	"go.uber.org/dig"
)

type UserServiceProvider struct{}

var _ _interface.ServiceProviderInterface = (*UserServiceProvider)(nil)

func (sp *UserServiceProvider) Register(container *dig.Container) error {
	if err := container.Provide(user_auth.NewAuthService, dig.Name("auth")); err != nil {
		return err
	}
	if err := container.Provide(user_service.NewUserService, dig.Name("user")); err != nil {
		return err
	}
	if err := container.Provide(user_cache.NewUserCache, dig.Name("user_cache")); err != nil {
		return err
	}
	if err := container.Provide(user_client.NewUserClient, dig.Name("user_client")); err != nil {
		return err
	}
	return nil
}
