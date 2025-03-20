// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 22:46

package user

import (
	"fmt"
	"github.com/duzhenlin/skittle/src/core/helper"
	jsoniter "github.com/json-iterator/go"
	"time"
)

func (u *User) UserAuth(id string) (interface{}, error) {
	return u.authRequest("id", id, "UserAuth")
}

func (u *User) UserAuthToken(token string) (interface{}, error) {
	return u.authRequest("token", token, "UserAuthToken")
}

// 通用认证请求方法
func (u *User) authRequest(paramKey, paramValue, method string) (interface{}, error) {
	// 构造请求参数
	args := map[string]interface{}{
		paramKey:    paramValue,
		"namespace": u.config.Skittle.Namespace,
	}

	// 序列化参数
	argsJSON, err := jsoniter.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("参数序列化失败: %w", err)
	}
	// 获取远程服务
	useService, err := u.hproseClient.GetDefaultUserService()
	if err != nil {
		return nil, fmt.Errorf("获取远程服务失败: %w", err)
	}
	// 执行远程调用
	var result interface{}

	switch method {
	case "UserAuth":
		result, err = useService.UserAuth(string(argsJSON), u.config.Skittle.ModuleId)
	case "UserAuthToken":
		result, err = useService.UserAuthToken(string(argsJSON), u.config.Skittle.ModuleId)
	default:
		return nil, fmt.Errorf("未知的认证方法: %s", method)
	}

	// 调试日志
	if u.config.Debug {
		helper.RunTime(time.Now(), u.hproseClient.TargetName, method, args)
		fmt.Printf("[DEBUG] 远程认证结果: %v\n", result)
	}

	// 错误处理
	if err != nil {
		return nil, fmt.Errorf("远程认证失败: %w", err)
	}

	// 处理响应数据
	return parseAuthResponse(result)
}

// 统一响应解析
func parseAuthResponse(result interface{}) (LoginRes, error) {
	responseData, err := jsoniter.Marshal(result)
	if err != nil {
		return LoginRes{}, fmt.Errorf("响应序列化失败: %w", err)
	}

	var authRes LoginRes
	if err := jsoniter.Unmarshal(responseData, &authRes); err != nil {
		return LoginRes{}, fmt.Errorf("响应解析失败: %w", err)
	}
	return authRes, nil
}
