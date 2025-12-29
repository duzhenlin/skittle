// Package examples
// 自定义选项示例：灵活配置模块

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/core"
)

func main() {
	ctx := context.Background()
	cfg := &config.Config{
		Debug: true,
		Redis: &config.Redis{
			ConType: "single",
			Addr:    "127.0.0.1:6379",
			Pwd:     "",
			Db:      0,
		},
		Skittle: &config.Skittle{
			Namespace: "custom-app",
		},
	}

	// 方式三：自定义选项
	fmt.Println("=== 方式三：自定义选项模式 ===")
	
	// 创建自定义选项：启用缓存和日志，禁用用户和Hprose
	opts := core.DefaultAppOptions().
		WithUserModule(false).      // 禁用用户模块
		WithHproseModule(false).     // 禁用Hprose
		WithCacheModule(true).       // 启用缓存
		WithLogModule(true)          // 启用日志

	app, err := core.NewAppWithOptions(ctx, cfg, opts)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	// 使用启用的模块
	app.Logger.Info(ctx, "自定义应用启动成功")
	
	// 使用缓存
	err = app.Cache.Set("custom", "value", 3600)
	if err == nil {
		fmt.Println("✓ 缓存功能正常")
	}

	// 检查模块状态
	fmt.Printf("用户模块: %v\n", app.HasUserModule())
	fmt.Printf("Hprose客户端: %v\n", app.Client != nil)
	fmt.Printf("缓存: %v\n", app.Cache != nil)
	fmt.Printf("日志: %v\n", app.Logger != nil)
}

