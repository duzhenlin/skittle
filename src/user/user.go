// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 18:37

package user

import (
	"context"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/redis"
	"time"
)

type User struct {
	config *config.Config
	IsLock int `json:"is_lock"`
}

var user *User

func init() {
	user = &User{}
}

// GetUserInstance 获取实例
func GetUserInstance(config *config.Config) *User {
	user.config = config
	return user
}

func (u *User) SetIsLock(isLock int) {
	u.IsLock = isLock
}

func (u *User) GetConfig() *config.Config {
	return u.config
}

func (u *User) SetUserInfo(UserInfo *LoginData) (string, error) {
	token := UserInfo.ModuleToken
	client := redis.GetClientInstance(u.config)
	ctx := context.Background()
	key := u.CacheKey("u_" + token)
	result, err := client.Set(ctx, key, UserInfo, 60*60*2*time.Second).Result()
	return result, err
}
func (u *User) CacheKey(key string) string {
	namespace := u.config.Skittle.Namespace
	return namespace + ":" + key
}
