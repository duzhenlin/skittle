// Package core
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 16:10

package core

import (
	"context"

	"github.com/duzhenlin/skittle/v2/src/config"
)

// NewAppWithDefaults 便捷方法：使用默认配置创建应用（即开即用）
// 这是最简单的使用方式，包含所有模块
func NewAppWithDefaults(ctx context.Context, cfg *config.Config) (*App, error) {
	container, err := BuildContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return NewApp(container)
}

// NewCoreAppOnly 便捷方法：创建仅核心功能的应用（不包含user模块）
// 适用于只需要缓存、日志等核心功能的场景
func NewCoreAppOnly(ctx context.Context, cfg *config.Config) (*App, error) {
	container, err := BuildContainerWithOptions(ctx, cfg, CoreOnlyOptions())
	if err != nil {
		return nil, err
	}
	return NewCoreApp(container)
}

// NewAppWithOptions 便捷方法：使用自定义选项创建应用
// 提供最大的灵活性
func NewAppWithOptions(ctx context.Context, cfg *config.Config, opts *AppOptions) (*App, error) {
	container, err := BuildContainerWithOptions(ctx, cfg, opts)
	if err != nil {
		return nil, err
	}
	
	// 根据是否启用user模块选择不同的创建方式
	if opts.EnableUserModule {
		return NewApp(container)
	}
	return NewCoreApp(container)
}

