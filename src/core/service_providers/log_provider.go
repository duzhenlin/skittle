// Package service_providers
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 11:19

package service_providers

import (
	_interface "github.com/duzhenlin/skittle/v2/src/core/service_providers/interface"
	"github.com/duzhenlin/skittle/v2/src/log"
	"go.uber.org/dig"
)

type LogServiceProvider struct{}

var _ _interface.ServiceProviderInterface = (*LogServiceProvider)(nil)

func (sp *LogServiceProvider) Register(container *dig.Container) error {
	if err := container.Provide(log.NewLogger, dig.Name("log")); err != nil {
		return err
	}
	return nil
}
