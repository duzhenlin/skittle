// Package core
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 17:32

package core

import (
	"fmt"
	"github.com/duzhenlin/skittle/src/config"

	"testing"
)

func TestApp(t *testing.T) {
	app := NewApp(&config.Config{
		Redis: config.Redis{
			Cont: "127.0.0.1:6379",
		},
		Skittle: config.Skittle{
			Namespace: "skittle",
			Server: config.Server{
				IsModule: false,
				Register: nil,
			},
			Client: []config.Client{
				{
					ClientName: "user",
					ClientType: "http",
					ClientUrl:  "http://fugui.sdy.qiludev.com/open/hprose/start",
				},
			},
			ModuleId:  "ef63713b81581a41c80789ae72cca3be",
			SecretKey: "87c8715c522954bd4beb3d25f4d52f2a",
		},
	})

	fmt.Println(app.User.GetConfig())
}
