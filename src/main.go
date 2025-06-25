// Package skittle
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:41

package main

import (
	"context"
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/core"
	"log"
	"net/http"
	"net/url"
	"runtime"
)

func main() {
	ctx := context.Background()
	app, err := core.NewApp(ctx, &config.Config{
		Debug: true,
		Redis: &config.Redis{
			ConType:      "sentinel",
			MasterName:   "T1",
			SentinelAddr: []string{"192.168.10.84:17001"},
			Pwd:          "Iqilu@11231",
			Db:           10,
		},
		EsConfig: &config.EsConfig{
			Address:   "http://192.168.10.77:9200",
			Username:  "elastic",
			Password:  "Iqilu@11231",
			IndexName: "fugui_log_v2",
			Channel:   "contract",
		},

		Skittle: &config.Skittle{
			Namespace: "cs",
			Server: &config.Server{
				IsModule: true,
				Register: &config.Register{
					Data: &[]config.RegisterData{
						{
							ID:   "1",
							Pid:  0,
							Name: "首页",
							Ops: []config.RegisterOps{
								{
									URL:    "https://fugui.sdy.qiludev.com/open/hprose/start",
									Method: 1,
								},
							},
						},
						{
							ID:   "2",
							Pid:  1,
							Name: "首页",
							Ops: []config.RegisterOps{
								{
									URL:    "https://fugui.sdy.qiludev.com/open/hprose/start",
									Method: 1,
								},
							},
						},
					},
					Content: &[]config.RegisterContent{
						{
							Type:    "3",
							Key:     "column_read",
							KeyName: "token",
						},
					},
					Quick: &[]config.RegisterQuick{
						{
							Name: "首页",
							Icon: "1",
							Path: "/hello",
						},
					},
				},
				UserUpdateFunc: func(id string, orgId string, action string) {
					fmt.Printf("更新用户信息: %s %s %s \n", id, orgId, action)
				},
			},
			Client: &[]config.Client{
				{
					ClientName: "user",
					ClientType: "https",
					ClientUrl:  "https://fugui.sdy.qiludev.com/open/hprose/start",
				},
			},
			ModuleId:  "ef63713b81581a41c80789ae72cca3be",
			SecretKey: "87c8715c522954bd",
		},
	})
	if err != nil {
		log.Fatalf("应用启动失败: %v", err)
	}

	//fmt.Println(app.User.GetConfig())
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		service := app.Server.Start()
		service.AddFunction("notice1", func() string {
			return "hello world"
		})
		service.ServeHTTP(w, r)

	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		// 获取查询参数
		query := r.URL.Query()
		// 获取单个参数的值
		str := query.Get("sign")

		enEscapeUrl, _ := url.QueryUnescape(str)

		sign, err := app.User.LoginSign(enEscapeUrl)

		fmt.Println("解密后的签名:", sign)
		if err != nil {
			fmt.Println("登陆验证失败")
			fmt.Println(w, err)
			return
		}
		app.Log.Log(r)

	})

	http.HandleFunc("/shouji", func(w http.ResponseWriter, r *http.Request) {

		// 获取单个参数的值
		sign, err := app.User.LoginSignByToken(r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Fprintln(w, sign.Token)

	})
	fmt.Printf("当前运行环境：%v \n", runtime.GOOS)
	fmt.Printf("运行环境：127.0.0.1:8813  \n")
	http.ListenAndServe("0.0.0.0:8813", nil)

}
