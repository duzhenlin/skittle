// Package log
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 15:58

package log

import (
	"context"
	"fmt"
	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/user/user_service"
	"go.uber.org/dig"
)

type LogDeps struct {
	dig.In
	Config *config.Config
	Ctx    context.Context
	User   *user_service.UserService `name:"user"`
}

func NewLogger(deps LogDeps) Logger {
	zapLogger := NewZapLogger()
	esLogger, err := NewESLogger(deps.Config)
	if err != nil {
		// 允许ES不可用时只用本地日志
		fmt.Printf("NewESLogger err: %s", err.Error())
		esLogger = nil
	}

	fmt.Printf("NewLogger: %s \n", deps.Config.EsConfig.Address)
	return NewDefaultLogger(deps, zapLogger, esLogger)
}
