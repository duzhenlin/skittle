// Package auth
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 13:28

package user_auth

import (
	"github.com/duzhenlin/skittle/v2/src/user/user_model"
	"net/http"
)

type IAuthService interface {
	LoginSign(sign string) (*user_model.LoginSign, error)
	LoginUserId(userID string, platform string) (*user_model.LoginSign, error)
	LoginSignByToken(r *http.Request) (*user_model.LoginSign, error)
}
