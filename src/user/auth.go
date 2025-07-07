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
	"github.com/duzhenlin/skittle/src/constant"
	"github.com/duzhenlin/skittle/src/core/helper"
	"github.com/duzhenlin/skittle/src/core/helper/aes"
	"github.com/hprose/hprose-golang/util"
	"net/http"
	"strings"
)

// LoginSign 登录签名
func (u *User) LoginSign(sign string) (*LoginSign, error) {
	// 解密签名信息
	decryptedSign, err := aes.Decrypt(sign, u.config.Skittle.SecretKey, constant.Base64URL)
	if err != nil {
		return nil, fmt.Errorf("签名解密失败: %w", err)
	}

	// 解析签名内容
	parts := strings.Split(decryptedSign, "|") // 最多分割成两部分
	if len(parts) < 1 || parts[0] == "" {
		return nil, errors.New("无效的签名格式")
	}

	// 获取并验证用户ID
	userID := strings.TrimSpace(parts[0])
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	fmt.Println("用户ID:", userID)
	// 执行登录流程
	userInfo, err := u.Login(userID, "web")

	fmt.Println("用户信息:", userInfo)
	if err != nil {
		return nil, err
	}
	token, err := u.LoginLogic(userInfo)
	if err != nil {
		return nil, fmt.Errorf("登录流程执行失败: %w", err)
	}

	return &LoginSign{Token: token}, nil
}

// LoginUserId 根据用户id进行登陆，获取token
func (u *User) LoginUserId(userID string, platform string) (*LoginSign, error) {

	userInfo, err := u.Login(userID, platform)
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return nil, fmt.Errorf("用户信息获取失败: ID=%s", userID)
	}
	// 执行登录逻辑
	token, err := u.LoginLogic(userInfo)
	if err != nil {
		return nil, fmt.Errorf("登录处理失败: %w", err)
	}

	// 构造返回结果
	return &LoginSign{
		Token: token,
	}, nil

}

// LoginSignByToken 根据token进行登陆，获取token
func (u *User) LoginSignByToken(r *http.Request) (*LoginSign, error) {
	// 从请求头获取令牌
	appToken := strings.TrimSpace(r.Header.Get("app_token"))
	if appToken == "" {
		return nil, errors.New("缺失认证令牌")
	}
	// 获取用户信息
	userInfo, err := u.Login(appToken, "app")
	if err != nil {
		return nil, err
	}

	if userInfo == nil {
		return nil, errors.New("用户信息获取失败")
	}

	// 执行登录逻辑
	newToken, err := u.LoginLogic(userInfo)
	if err != nil {
		return nil, fmt.Errorf("令牌生成失败: %w", err)
	}

	// 返回登录凭证
	return &LoginSign{
		Token: newToken,
	}, nil
}

func (u *User) Login(id string, platform string) (interface{}, error) {

	// 验证输入参数
	if id == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 定义平台处理策略
	type authFunc func(string) (interface{}, error)
	strategies := map[string]authFunc{
		"web": u.UserAuth,
		"app": u.UserAuthToken,
		// 可扩展其他平台："MiniProgram",
	}

	// 获取对应的认证方法
	auth, ok := strategies[platform]
	if !ok {
		return nil, fmt.Errorf("不支持的平台类型: %s", platform)
	}

	// 执行认证逻辑
	info, err := auth(id)
	if err != nil {

		return nil, fmt.Errorf("认证服务不可用: %w", err)
	}

	return info, nil
}

func (u *User) LoginLogic(userInfo interface{}) (string, error) {

	// 类型安全检查
	res, ok := userInfo.(LoginRes)
	if !ok || !checkIsUserAuthRes(userInfo) {
		return "", fmt.Errorf("无效的用户认证类型: %T", userInfo)
	}

	// 获取用户数据
	UserInfo := res.LoginData

	// 存储用户信息
	if _, err := u.SetUserInfo(&UserInfo); err != nil {
		return "", fmt.Errorf("用户信息存储失败: %w", err)
	}

	// 生成或返回现有令牌
	if UserInfo.ModuleToken == "" {
		return generateModuleToken(UserInfo.ID, u.config.Skittle.Namespace), nil
	}

	return UserInfo.ModuleToken, nil
}

// generateModuleToken 生成模块访问令牌
func generateModuleToken(userID, namespace string) string {
	builder := strings.Builder{}
	builder.WriteString(util.UUIDv4())
	builder.WriteString(userID)
	builder.WriteString(namespace)
	return helper.GetStringMd5(builder.String())
}

func checkIsUserAuthRes(data interface{}) bool {
	switch data.(type) {
	case LoginRes:
		return true
	default:
		return false
	}
}
