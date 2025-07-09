// Package service
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 13:28

package user_service

import (
	"github.com/duzhenlin/skittle/v2/src/user/user_model"
	"net/http"
)

type IUserCache interface {
	SetUserInfo(user *user_model.LoginData) (string, error)
	GetUserInfo(token string) (interface{}, error)
	GetUserInfoById(userId string) (interface{}, error)
	GetUserInfoNew(h *http.Request) (user_model.LoginData, error)
}
