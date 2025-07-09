// Package service_providers
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 10:57

package core

import (
	"github.com/duzhenlin/skittle/v2/src/core/service_providers/interface"
	"go.uber.org/dig"
)

var Providers []_interface.ServiceProviderInterface

func RegisterProvider(provider _interface.ServiceProviderInterface) {
	Providers = append(Providers, provider)
}

func RegisterAll(container *dig.Container) error {
	for _, provider := range Providers {
		if err := provider.Register(container); err != nil {
			return err
		}
	}
	return nil
}
