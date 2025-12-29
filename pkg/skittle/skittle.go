// Package skittle 提供统一的导出入口
// 这是推荐的使用方式，提供简洁的 API
//
// 使用示例：
//   import "github.com/duzhenlin/skittle/v2/pkg/skittle"
//
//   ctx := context.Background()
//   cfg := &config.Config{...}
//   app, err := skittle.New(ctx, cfg)
package skittle

import (
	"context"

	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/core"
)

// New 创建应用实例（即开即用模式）
// 这是最简单的使用方式，包含所有模块（cache, log, user等）
//
// 示例：
//   ctx := context.Background()
//   cfg := &config.Config{...}
//   app, err := skittle.New(ctx, cfg)
//   if err != nil {
//       log.Fatal(err)
//   }
func New(ctx context.Context, cfg *config.Config) (*core.App, error) {
	return core.NewAppWithDefaults(ctx, cfg)
}

// NewCore 创建仅核心功能的应用（不包含user模块）
// 适用于只需要缓存、日志等核心功能的场景
//
// 示例：
//   ctx := context.Background()
//   cfg := &config.Config{...}
//   app, err := skittle.NewCore(ctx, cfg)
//   if err != nil {
//       log.Fatal(err)
//   }
func NewCore(ctx context.Context, cfg *config.Config) (*core.App, error) {
	return core.NewCoreAppOnly(ctx, cfg)
}

// NewWithOptions 使用自定义选项创建应用
// 提供最大的灵活性，可以按需配置模块
//
// 示例：
//   opts := core.DefaultAppOptions().
//       WithUserModule(false).
//       WithHproseModule(true)
//   app, err := skittle.NewWithOptions(ctx, cfg, opts)
func NewWithOptions(ctx context.Context, cfg *config.Config, opts *core.AppOptions) (*core.App, error) {
	return core.NewAppWithOptions(ctx, cfg, opts)
}

// App 是应用实例的别名，方便使用
type App = core.App

// AppOptions 是应用选项的别名，方便使用
type AppOptions = core.AppOptions

// DefaultAppOptions 返回默认选项（即开即用模式）
func DefaultAppOptions() *AppOptions {
	return core.DefaultAppOptions()
}

// CoreOnlyOptions 返回仅核心模块选项（不包含业务逻辑）
func CoreOnlyOptions() *AppOptions {
	return core.CoreOnlyOptions()
}

