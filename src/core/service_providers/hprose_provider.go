// Package service_providers
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 10:55

package service_providers

import (
	"github.com/duzhenlin/skittle/v2/src/core/service_providers/interface"
	"github.com/duzhenlin/skittle/v2/src/hprose/hprose_client"
	"github.com/duzhenlin/skittle/v2/src/hprose/hprose_server"
	"go.uber.org/dig"
)

type HproseServiceProvider struct{}

var _ _interface.ServiceProviderInterface = (*HproseServiceProvider)(nil)

func (sp *HproseServiceProvider) Register(container *dig.Container) error {
	if err := container.Provide(hprose_client.NewClient, dig.Name("hprose_client")); err != nil {
		return err
	}
	if err := container.Provide(hprose_server.NewServer, dig.Name("hprose_server")); err != nil {
		return err
	}
	return nil
}
