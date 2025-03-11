// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 22:46

package user

import (
	"github.com/duzhenlin/skittle/src/core/helper"
	"github.com/duzhenlin/skittle/src/hprose/client"
	jsoniter "github.com/json-iterator/go"
	"time"
)

func (u *User) UserAuth(id string) (interface{}, error) {

	httpClient := &client.Client{}
	httpClient.SetConfig(u.config)

	args := make(map[string]interface{})
	args["id"] = id
	args["namespace"] = u.config.Skittle.Namespace
	argsJson, _ := jsoniter.Marshal(args)
	defer helper.RunTime(time.Now(), httpClient.ClientName, "Index", args)

	result, err := httpClient.GetDefaultUserService().UserAuth(string(argsJson), u.config.Skittle.ModuleId)
	if err != nil {
		return nil, err
	}
	arr, err := jsoniter.Marshal(result)
	if err != nil {
		return nil, err
	}
	var UserAuthRes LoginRes
	err = jsoniter.Unmarshal(arr, &UserAuthRes)
	if err != nil {
		return nil, err
	}
	return UserAuthRes, err
}

func (u *User) UserAuthToken(token string) (interface{}, error) {

	httpClient := &client.Client{}
	httpClient.SetConfig(u.config)
	httpClient.To("user")

	args := make(map[string]interface{})
	args["token"] = token
	args["namespace"] = u.config.Skittle.Namespace
	argsJson, _ := jsoniter.Marshal(args)
	defer helper.RunTime(time.Now(), httpClient.ClientName, "Index", args)

	result, err := httpClient.GetDefaultUserService().UserAuthToken(string(argsJson), u.config.Skittle.ModuleId)
	if err != nil {
		return nil, err
	}
	arr, err := jsoniter.Marshal(result)
	if err != nil {
		return nil, err
	}
	var UserAuthRes LoginRes
	err = jsoniter.Unmarshal(arr, &UserAuthRes)
	if err != nil {
		return nil, err
	}
	return UserAuthRes, err
}
