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
	if token == "" {
		log.Println("ExtendTime token为空")
		return false
	}
	key := c.CacheKey("u_" + token)
	result, err := c.cache.Get(key)
	if err != nil {
		log.Println("ExtendTime 获取缓存失败:", err)
		return false
	}
	if result == "" {
		log.Println("ExtendTime 缓存不存在")
		return false
	}
	// 只续期，不重复写入
	if err := c.cache.Expire(key, int((2 * time.Hour).Seconds())); err != nil {
		log.Println("ExtendTime 续期失败:", err)
		return false
	}
	return true
}

func (c *UserService) GetUserInfo(token string) (interface{}, error) {
	if c.Ctx == nil {
		fmt.Println("上下文未初始化")
	}
	key := c.CacheKey("u_" + token)
	helper.DebugLog(c.config, "[UserService] 获取缓存key：%s", key)
	result, err := c.cache.Get(key)
	helper.DebugLog(c.config, "[UserService] 获取缓存结果: %v", result)
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

func (c *UserService) GetUserInfoNew(r *http.Request) (user_model.LoginData, error) {
	var resultStr string

	jwtToken := r.Header.Get("Jwttoken")
	if jwtToken != "" {
		userId, err := c.DeCodeJwtToken(jwtToken)
		if err != nil {
			return user_model.LoginData{}, fmt.Errorf("JWT解析失败: %w", err)
		}
		if userId == "" {
			return user_model.LoginData{}, fmt.Errorf("无效的用户ID")
		}
		// 从缓存中获取用户信息
		res, err := c.GetUserInfoById(userId)
		if err != nil {
			return user_model.LoginData{}, fmt.Errorf("获取用户信息失败: %w", err)
		}
		var ok bool
		resultStr, ok = res.(string)
		if !ok {
			return user_model.LoginData{}, fmt.Errorf("缓存数据类型错误")
		}
	} else {
		token := r.Header.Get("token")
		if token == "" {
			return user_model.LoginData{}, fmt.Errorf("未提供 token")
		}
		helper.DebugLog(c.config, "[UserService] 执行GetUserInfoNew: GetUserInfo")
		res, err := c.GetUserInfo(token)
		if err != nil {
			return user_model.LoginData{}, fmt.Errorf("获取用户信息失败: %w", err)
		}
		var ok bool
		resultStr, ok = res.(string)
		if !ok {
			return user_model.LoginData{}, fmt.Errorf("缓存数据类型错误")
		}
	}

	helper.DebugLog(c.config, "[UserService] 获取的用户信息: %v", resultStr)

	// 解析用户信息
	var userData user_model.LoginData
	if err := jsoniter.Unmarshal([]byte(resultStr), &userData); err != nil {
		return user_model.LoginData{}, fmt.Errorf("解析用户信息失败: %w", err)
	}

	return userData, nil
}

// DeCodeJwtToken 解析JWT Token并返回用户ID
// 从配置中读取JWT密钥，支持密钥配置化
func (c *UserService) DeCodeJwtToken(jwtToken string) (string, error) {
	if jwtToken == "" {
		return "", fmt.Errorf("JWT token为空")
	}

	// 从配置读取JWT密钥
	jwtSecret := c.config.Skittle.JwtSecret
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT密钥未配置")
	}

	// 解析JWT Token
	token, err := jwt.ParseWithClaims(jwtToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("JWT解析失败: %w", err)
	}

	// 安全类型断言
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("无效的JWT claims")
	}

	// 提取用户ID
	userID, ok := claims["userid"]
	if !ok {
		return "", fmt.Errorf("JWT claims中缺少userid字段")
	}

	return helper.GetInterfaceToString(userID), nil
}

func (c *UserService) CacheKey(suffix string) string {
	return fmt.Sprintf("%s:%s", c.config.Skittle.Namespace, suffix)
}
