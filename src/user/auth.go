// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:39

package user

import (
	"errors"
	"fmt"
	"github.com/duzhenlin/skittle/src/core/helper"
	"github.com/hprose/hprose-golang/util"
	"net/http"
	"strings"
)

func (u *User) LoginSign(sign string) (*LoginSign, error) {
	signEncryptStr, err := helper.AesDecrypt(sign, u.config.Skittle.SecretKey)
	if err != nil {
		return nil, errors.New("登录失败，密钥解密失败")
	}

	strArrayNew := strings.Split(signEncryptStr, "|")
	if len(strArrayNew) >= 1 {
		//获取用户id
		id := strArrayNew[0]
		userInfo := u.Login(id, "web")
		//登陆逻辑
		token, err := u.LoginLogic(userInfo)

		if err != nil {
			return nil, err
		}
		ret := new(LoginSign)
		ret.Token = token
		return ret, nil
	}
	return nil, errors.New("登录失败，密钥解密失败")
}

func (u *User) LoginSignByToken(r *http.Request) (*LoginSign, error) {
	token := r.Header.Get("app_token")
	userInfo := u.Login(token, "app")
	token, err := u.LoginLogic(userInfo)

	if err != nil {
		return nil, err
	}
	ret := new(LoginSign)
	ret.Token = token
	return ret, nil
}

func (u *User) Login(id string, platform string) interface{} {
	//namespace := u.config.Skittle.Namespace
	var info interface{}
	if platform == "web" {
		info, err := u.UserAuth(id)
		if err != nil {
			fmt.Println(err)
		}
		return info
	}
	//else if platform == "MiniProgram" {
	//	info, err := PcUserAuth(id)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	return info
	if platform == "app" {
		info, err := u.UserAuthToken(id)
		if err != nil {
			fmt.Println(err)
		}
		return info
	}

	return info
}

func (u *User) LoginLogic(userInfo interface{}) (string, error) {

	if !u.checkIsUserAuthRes(userInfo) {
		return "", errors.New("登陆失败")
	}
	UserInfo := userInfo.(LoginRes).LoginData
	_, err := u.SetUserInfo(&UserInfo)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("登陆失败")
	}
	if UserInfo.ModuleToken == "" {
		srt := fmt.Sprintf("%s%s%s", util.UUIDv4(), UserInfo.ID, u.config.Skittle.Namespace)
		return helper.GetStringMd5(srt), nil
	}
	return UserInfo.ModuleToken, nil
}

func (u *User) checkIsUserAuthRes(data interface{}) bool {
	switch data.(type) {
	case LoginRes:
		return true
	default:
		return false
	}
}
