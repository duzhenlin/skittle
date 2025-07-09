// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 22:46

package user_client

import (
	"context"
	"fmt"
	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/core/helper"
	"github.com/duzhenlin/skittle/v2/src/hprose/hprose_client"
	"github.com/duzhenlin/skittle/v2/src/user/user_model"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/dig"
	"time"
)

type ServiceDefaultUser struct {
	UserAuth                 func(string, string) (interface{}, error) `name:"user_auth"`
	UserAuthToken            func(string, string) (interface{}, error) `name:"user_auth_token"`
	UserAuthMiniProgramToken func(string, string) (interface{}, error) `name:"user_auth_miniprogram_token"`
}

type UserClientService struct {
	Ctx          context.Context
	hproseClient *hprose_client.HproseClientService
	config       *config.Config // 你的 config 类型
}

type ClientDeps struct {
	dig.In
	Ctx          context.Context
	HproseClient *hprose_client.HproseClientService `name:"hprose_client"`
	Config       *config.Config
}

func NewUserClient(deps ClientDeps) *UserClientService {
	return &UserClientService{
		Ctx:          deps.Ctx,
		hproseClient: deps.HproseClient,
		config:       deps.Config,
	}
}

func (c *UserClientService) getDefaultUserService() (ServiceDefaultUser, error) {
	var serviceDefaultUser ServiceDefaultUser
	c.hproseClient.WithTarget("user")
	serviceClient, err := c.hproseClient.GetClient()
	if err != nil {
		return serviceDefaultUser, err
	}
	serviceClient.UseService(&serviceDefaultUser)
	return serviceDefaultUser, nil
}

func (c *UserClientService) authRequest(paramKey, paramValue, method string) (interface{}, error) {
	// 构造请求参数
	args := map[string]interface{}{
		paramKey:    paramValue,
		"namespace": c.config.Skittle.Namespace,
	}
	// 序列化参数
	argsJSON, err := jsoniter.Marshal(args)
	if err != nil {
		return user_model.LoginRes{}, fmt.Errorf("参数序列化失败: %w", err)
	}
	// 获取远程服务
	useService, err := c.getDefaultUserService()
	if err != nil {
		return nil, fmt.Errorf("获取远程服务失败: %w", err)
	}

	var result interface{}
	switch method {
	case "UserAuth":
		result, err = useService.UserAuth(string(argsJSON), c.config.Skittle.ModuleId)
	case "UserAuthToken":
		result, err = useService.UserAuthToken(string(argsJSON), c.config.Skittle.ModuleId)
	default:
		return nil, fmt.Errorf("未知的认证方法: %s", method)
	}
	// 调试日志
	if c.config.Debug {
		helper.RunTime(time.Now(), c.hproseClient.TargetName, method, args)
		fmt.Printf("[DEBUG] 远程认证结果: %v\n", result)
	}
	if err != nil {
		return nil, fmt.Errorf("远程认证失败: %w", err)
	}
	return parseAuthResponse(result)
}

func (c *UserClientService) UserAuth(id string) (interface{}, error) {
	return c.authRequest("id", id, "UserAuth")
}

func (c *UserClientService) UserAuthToken(token string) (interface{}, error) {
	return c.authRequest("token", token, "UserAuthToken")
}

func (c *UserClientService) UserAuthMiniProgramToken(token string) (interface{}, error) {
	return c.authRequest("token", token, "UserAuthMiniProgramToken")
}

func parseAuthResponse(result interface{}) (user_model.LoginRes, error) {
	responseData, err := jsoniter.Marshal(result)
	if err != nil {
		return user_model.LoginRes{}, fmt.Errorf("响应序列化失败: %w", err)
	}
	var authRes user_model.LoginRes
	if err := jsoniter.Unmarshal(responseData, &authRes); err != nil {
		return user_model.LoginRes{}, fmt.Errorf("响应解析失败: %w", err)
	}
	return authRes, nil
}
