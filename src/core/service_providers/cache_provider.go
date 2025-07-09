// Package service_providers
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 11:26

package service_providers

import (
	"github.com/duzhenlin/skittle/v2/src/cache"
	_interface "github.com/duzhenlin/skittle/v2/src/core/service_providers/interface"
	"go.uber.org/dig"
)

type CacheServiceProvider struct{}

var _ _interface.ServiceProviderInterface = (*CacheServiceProvider)(nil)

func (sp *CacheServiceProvider) Register(container *dig.Container) error {
	if err := container.Provide(cache.NewCache, dig.Name("cache")); err != nil {
		return err
	}

	return nil
}
