// Package core
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:17

package core

import (
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/hprose"
	"github.com/duzhenlin/skittle/src/hprose/client"
	"github.com/duzhenlin/skittle/src/user"
	"strings"
)

const providerList = "config,Client,Server,User"

type App struct {
	Config *config.Config
	Client *client.Client
	Server *hprose.Server
	User   *user.User
}

func NewApp(config *config.Config) *App {
	a := &App{}
	a.Config = config
	a.Config.SetConfig(*config)

	err := a.RegisterProviders()
	if err != nil {
		return nil
	}
	return a
}

func (a *App) RegisterProviders() error {

	provider := strings.Split(providerList, ",")
	for _, p := range provider {
		switch p {
		case "Client":
			a.Client = client.GetClientInstance(a.Config)
		case "Server":
			a.Server = hprose.GetServerInstance(a.Config)
		case "User":
			a.User = user.GetUserInstance(a.Config)
		}
		fmt.Println(p)
	}
	return nil
}
