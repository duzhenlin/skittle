// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 18:37

package user

import (
	"context"
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/core/helper"
	"github.com/duzhenlin/skittle/src/hprose/client"
	"github.com/forgoer/openssl"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"time"
)

// User 用户服务结构体
type User struct {
	ctx          context.Context
	config       *config.Config
	redisClient  *redis.Client  // 注入的Redis客户端
	hproseClient *client.Client // Hprose客户端
}

// NewUser 创建用户服务实例
func NewUser(ctx context.Context, cfg *config.Config) *User {
	return &User{
		ctx:    ctx,
		config: cfg,
	}
}

// SetHproseClient 设置hprose客户端
func (u *User) SetHproseClient(client *client.Client) *User {
	u.hproseClient = client
	return u
}

// SetRedisClient 设置redis客户端
func (u *User) SetRedisClient(redisClient *redis.Client) *User {
	u.redisClient = redisClient
	return u
}

func (u *User) GetConfig() *config.Config {
	return u.config
}

// CacheKey 生成缓存键
func (u *User) CacheKey(suffix string) string {
	return fmt.Sprintf("%s:%s", u.config.Skittle.Namespace, suffix)
}

func (u *User) SetUserInfo(userInfo *LoginData) (string, error) {
	fmt.Println("设置缓存key：" + userInfo.ModuleToken)
	// 设置缓存
	key := u.CacheKey("u_" + userInfo.ModuleToken)
	result, err := u.redisClient.Set(u.ctx, key, userInfo, 2*time.Hour).Result()
	if err != nil {
		return "", fmt.Errorf("用户缓存写入失败: %w", err)
	}
	return result, err
}

func (u *User) GetUserInfo(token string) (interface{}, error) {
	if u.redisClient == nil {
		fmt.Println("Redis客户端未初始化")
	}
	if u.ctx == nil {
		fmt.Println("上下文未初始化")
	}
	key := u.CacheKey("u_" + token)
	fmt.Println("获取缓存key：" + key)
	result, err := u.redisClient.Get(u.ctx, key).Result()
	fmt.Println("获取缓存结果:", result)
	if err != nil {
		return nil, fmt.Errorf("用户缓存读取失败: %w", err)
	}
	return result, err
}
func (u *User) GetUserInfoById(userId string) (interface{}, error) {
	key := u.CacheKey("u_" + userId)
	result, err := u.redisClient.Get(u.ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("用户缓存读取失败: %w", err)
	}
	return result, err
}

func (u *User) GetUserInfoNew(c *http.Request) (LoginData, error) {
	var err error
	var result interface{}

	jwtToken := c.Header.Get("Jwttoken")

	if jwtToken != "" {
		userId := DeCodeJwtToken(jwtToken)
		if userId == "" {
			return LoginData{}, fmt.Errorf("无效的用户ID")
		}
		// 从缓存中获取用户信息
		result, err = u.GetUserInfoById(userId)

	} else {
		token := c.Header.Get("token")
		result, err = u.GetUserInfo(token)
	}
	fmt.Println("获取的用户信息:", result)
	if err != nil {
		return LoginData{}, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 解析用户信息
	var userData LoginData
	err = jsoniter.Unmarshal([]byte(result.(string)), &userData)
	if err != nil {
		return LoginData{}, fmt.Errorf("解析用户信息失败: %w", err)
	}

	return userData, nil
}

func GetJwtToken(userid string) (string, error) {
	var hmacSampleSecret []byte
	now := carbon.Now(carbon.Shanghai)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "jwt_admin",
		//主题；
		"sub": "jwt_admin",
		//接收方；
		"aud":    "all",
		"exp":    now.Timestamp() + 60*10,
		"nbf":    now.Timestamp(),
		"iat":    now.Timestamp(),
		"jti":    openssl.Md5("555"),
		"userid": userid,
	})

	hmacSampleSecret = []byte("iqilu@fugui")
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

func DeCodeJwtToken(jwtToken string) string {
	token, err := jwt.ParseWithClaims(jwtToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte("iqilu@fugui"), nil })
	if err != nil {
		return ""
	}
	return helper.GetInterfaceToString(token.Claims.(jwt.MapClaims)["userid"])
}
