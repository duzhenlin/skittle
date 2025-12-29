// Package examples
// 基础使用示例：展示即开即用模式

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
			Namespace: "example-app",
			SecretKey: "your-secret-key",
			JwtSecret: "your-jwt-secret",
		},
	}

	// 方式一：即开即用（推荐）- 包含所有模块
	fmt.Println("=== 方式一：即开即用模式 ===")
	app, err := core.NewAppWithDefaults(ctx, cfg)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	// 使用日志
	app.Logger.Info(ctx, "应用启动成功")

	// 使用缓存
	err = app.Cache.Set("hello", "world", 3600)
	if err != nil {
		log.Printf("设置缓存失败: %v", err)
	}

	value, err := app.Cache.Get("hello")
	if err == nil {
		fmt.Printf("缓存值: %s\n", value)
	}

	// 检查是否启用了用户模块
	if app.HasUserModule() {
		fmt.Println("✓ 用户模块已启用")
		// 可以使用 app.Auth, app.User 等
	} else {
		fmt.Println("✗ 用户模块未启用")
	}
}

