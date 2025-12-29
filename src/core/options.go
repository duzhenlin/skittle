// Package core
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 16:00

package core

// AppOptions 应用初始化选项
type AppOptions struct {
	// EnableUserModule 是否启用用户模块（默认true，保持即开即用）
	EnableUserModule bool
	// EnableHproseModule 是否启用Hprose模块（默认true）
	EnableHproseModule bool
	// EnableLogModule 是否启用日志模块（默认true）
	EnableLogModule bool
	// EnableCacheModule 是否启用缓存模块（默认true）
	EnableCacheModule bool
}

// DefaultAppOptions 返回默认选项（即开即用模式）
func DefaultAppOptions() *AppOptions {
	return &AppOptions{
		EnableUserModule:   true,
		EnableHproseModule: true,
		EnableLogModule:    true,
		EnableCacheModule:  true,
	}
}

// CoreOnlyOptions 返回仅核心模块选项（不包含业务逻辑）
func CoreOnlyOptions() *AppOptions {
	return &AppOptions{
		EnableUserModule:   false,
		EnableHproseModule: true,
		EnableLogModule:    true,
		EnableCacheModule:  true,
	}
}

// WithUserModule 启用用户模块
func (o *AppOptions) WithUserModule(enable bool) *AppOptions {
	o.EnableUserModule = enable
	return o
}

// WithHproseModule 启用Hprose模块
func (o *AppOptions) WithHproseModule(enable bool) *AppOptions {
	o.EnableHproseModule = enable
	return o
}

// WithLogModule 启用日志模块
func (o *AppOptions) WithLogModule(enable bool) *AppOptions {
	o.EnableLogModule = enable
	return o
}

// WithCacheModule 启用缓存模块
func (o *AppOptions) WithCacheModule(enable bool) *AppOptions {
	o.EnableCacheModule = enable
	return o
}

