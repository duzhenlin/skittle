// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 18:37

package user_service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/duzhenlin/skittle/v2/src/cache"
	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/core/helper"
	"github.com/duzhenlin/skittle/v2/src/user/user_client"
	"github.com/duzhenlin/skittle/v2/src/user/user_model"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/dig"
)

type UserService struct {
	Ctx    context.Context
	client *user_client.UserClientService
	cache  cache.Cache
	config *config.Config
}

type UserServiceDeps struct {
	dig.In
	Client *user_client.UserClientService `name:"user_client"`
	Ctx    context.Context
	Config *config.Config
	Cache  cache.Cache `name:"cache"`
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		Ctx:    deps.Ctx,
		client: deps.Client,
		config: deps.Config,
		cache:  deps.Cache,
	}
}

func (c *UserService) SetUserInfo(userInfo *user_model.LoginData) (string, error) {
	key := c.CacheKey("u_" + userInfo.ModuleToken)
	fmt.Printf("设置缓存key: %s \n", key)
	data, err := jsoniter.Marshal(userInfo)
	if err != nil {
		return "", fmt.Errorf("用户缓存写入失败: %w \n", err)
	}
	err = c.cache.Set(key, string(data), int((2 * time.Hour).Seconds()))
	return key, err
}
func (c *UserService) ExtendTime(token string) bool {
	if token != "" {
		key := c.CacheKey("u_" + token)
		result, err := c.cache.Get(key)
		if err != nil {
			log.Println("ExtendTime 获取缓存失败:", err)
			return false
		}
		if result != "" {
			// 解析用户信息
			var userData user_model.LoginData
			err = jsoniter.Unmarshal([]byte(result), &userData)
			if err != nil {
				log.Println("ExtendTime 解析用户信息失败:", err)
				return false
			}
			// 存储用户信息
			if _, err = c.SetUserInfo(&userData); err != nil {
				log.Println("ExtendTime 存储用户信息失败:", err)
				return false
			}
		}
		return true
	} else {
		log.Println("ExtendTime token为空")
		return false
	}
}

func (c *UserService) GetUserInfo(token string) (interface{}, error) {
	if c.Ctx == nil {
		fmt.Println("上下文未初始化")
	}
	key := c.CacheKey("u_" + token)
	fmt.Println("获取缓存key：" + key)
	result, err := c.cache.Get(key)
	fmt.Println("获取缓存结果:", result)
	if err != nil {
		return nil, fmt.Errorf("用户缓存读取失败: %w", err)
	}
	return result, err
}

func (c *UserService) GetUserInfoById(userId string) (interface{}, error) {
	key := c.CacheKey("u_" + userId)
	result, err := c.cache.Get(key)
	if err != nil {
		return nil, fmt.Errorf("用户缓存读取失败: %w", err)
	}
	return result, err
}

func (c *UserService) GetUserInfoNew(h *http.Request) (user_model.LoginData, error) {
	var err error
	var result interface{}

	jwtToken := h.Header.Get("Jwttoken")

	if jwtToken != "" {
		userId := DeCodeJwtToken(jwtToken)
		if userId == "" {
			return user_model.LoginData{}, fmt.Errorf("无效的用户ID")
		}
		// 从缓存中获取用户信息
		result, err = c.GetUserInfoById(userId)

	} else {
		token := h.Header.Get("token")
		result, err = c.GetUserInfo(token)
	}
	fmt.Println("获取的用户信息:", result)
	if err != nil {
		return user_model.LoginData{}, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 解析用户信息
	var userData user_model.LoginData
	err = jsoniter.Unmarshal([]byte(result.(string)), &userData)
	if err != nil {
		return user_model.LoginData{}, fmt.Errorf("解析用户信息失败: %w", err)
	}

	return userData, nil
}

func DeCodeJwtToken(jwtToken string) string {
	token, err := jwt.ParseWithClaims(jwtToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte("iqilu@fugui"), nil })
	if err != nil {
		return ""
	}
	return helper.GetInterfaceToString(token.Claims.(jwt.MapClaims)["userid"])
}

func (c *UserService) CacheKey(suffix string) string {
	return fmt.Sprintf("%s:%s", c.config.Skittle.Namespace, suffix)
}
