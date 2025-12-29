// Package examples
// 核心功能示例：仅使用核心模块，不包含user业务逻辑

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
	}

	// 方式二：仅核心功能（不包含user模块）
	fmt.Println("=== 方式二：仅核心功能模式 ===")
	app, err := core.NewCoreAppOnly(ctx, cfg)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	// 使用核心功能
	app.Logger.Info(ctx, "核心应用启动成功")

	// 使用缓存
	err = app.Cache.Set("key", "value", 3600)
	if err != nil {
		log.Printf("设置缓存失败: %v", err)
	}

	// 检查用户模块
	if app.HasUserModule() {
		fmt.Println("✗ 用户模块不应该启用")
	} else {
		fmt.Println("✓ 用户模块已禁用（符合预期）")
	}
}

