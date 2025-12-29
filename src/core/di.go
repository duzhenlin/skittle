// Package core
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/7
// @Time: 18:55

package core

import (
	"context"

	"github.com/duzhenlin/skittle/v2/src/config"
	"go.uber.org/dig"
)

// BuildContainer 构建容器（使用默认选项，即开即用模式）
func BuildContainer(ctx context.Context, cfg *config.Config) (*dig.Container, error) {
	return BuildContainerWithOptions(ctx, cfg, DefaultAppOptions())
}

// BuildContainerWithOptions 构建容器（支持自定义选项）
func BuildContainerWithOptions(ctx context.Context, cfg *config.Config, opts *AppOptions) (*dig.Container, error) {
	c := dig.New()

	// 基础依赖
	if err := c.Provide(func() context.Context { return ctx }); err != nil {
		return nil, err
	}
	if err := c.Provide(func() *config.Config { return cfg }); err != nil {
		return nil, err
	}

	// 根据选项注册服务
	if err := RegisterWithOptions(c, opts); err != nil {
		return nil, err
	}

	return c, nil
}
