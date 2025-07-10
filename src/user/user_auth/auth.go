// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:39

package user_auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/constant"
	"github.com/duzhenlin/skittle/v2/src/core/helper"
	"github.com/duzhenlin/skittle/v2/src/core/helper/aes"
	"github.com/duzhenlin/skittle/v2/src/user/user_client"
	"github.com/duzhenlin/skittle/v2/src/user/user_model"
	"github.com/duzhenlin/skittle/v2/src/user/user_service"
	"go.uber.org/dig"
)

// 平台类型常量，避免魔法字符串
const (
	PlatformWeb         = "web"
	PlatformApp         = "app"
	PlatformMiniProgram = "MiniProgram"
)

// AuthService 认证服务
// 负责用户认证、登录、Token生成等逻辑
// 依赖注入建议通过 AuthDeps 结构体
type AuthService struct {
	config     *config.Config
	user       *user_service.UserService
	userClient *user_client.UserClientService
	Ctx        context.Context
}

type AuthDeps struct {
	dig.In
	Ctx        context.Context
	Cfg        *config.Config
	User       *user_service.UserService      `name:"user"`
	UserClient *user_client.UserClientService `name:"user_client"`
}

// NewAuthService 构造函数
func NewAuthService(deps AuthDeps) *AuthService {
	return &AuthService{
		config:     deps.Cfg,
		user:       deps.User,
		userClient: deps.UserClient,
		Ctx:        deps.Ctx,
	}
}

// LoginSign 登录签名，解密签名并登录
func (a *AuthService) LoginSign(sign string) (*user_model.LoginSign, error) {
	decrypted, err := aes.Decrypt(sign, a.config.Skittle.SecretKey, constant.Base64URL)
	if err != nil {
		return nil, fmt.Errorf("签名解密失败: %w", err)
	}
	parts := strings.SplitN(decrypted, "|", 2)
	userID := strings.TrimSpace(parts[0])
	helper.DebugLog(a.config, "[AuthService] LoginSign 解密后用户ID: %s", userID)
	if userID == "" {
		return nil, errors.New("无效的签名格式或用户ID为空")
	}
	return a.LoginUserId(userID, PlatformWeb)
}

// LoginUserId 根据用户id进行登陆，获取token
func (a *AuthService) LoginUserId(userID, platform string) (*user_model.LoginSign, error) {
	helper.DebugLog(a.config, "[AuthService] LoginUserId userID: %s, platform: %s", userID, platform)
	userInfo, err := a.Login(userID, platform)
	if err != nil {
		helper.DebugLog(a.config, "[AuthService] LoginUserId 登录失败: %v", err)
		return nil, err
	}
	helper.DebugLog(a.config, "[AuthService] LoginUserId userInfo: %+v", userInfo)
	token, err := a.LoginLogic(userInfo)
	if err != nil {
		helper.DebugLog(a.config, "[AuthService] LoginUserId 登录逻辑失败: %v", err)
		return nil, fmt.Errorf("登录处理失败: %w", err)
	}
	helper.DebugLog(a.config, "[AuthService] LoginUserId token: %s", token)
	return &user_model.LoginSign{Token: token}, nil
}

// LoginSignByToken 根据token进行登陆，获取token
func (a *AuthService) LoginSignByToken(r *http.Request) (*user_model.LoginSign, error) {
	appToken := strings.TrimSpace(r.Header.Get("appToken"))
	helper.DebugLog(a.config, "[AuthService] LoginSignByToken appToken: %s", appToken)
	if appToken == "" {
		return nil, errors.New("缺失认证令牌")
	}
	return a.LoginUserId(appToken, PlatformApp)
}

// Login 认证入口，根据平台类型选择认证方式
func (a *AuthService) Login(id, platform string) (interface{}, error) {
	helper.DebugLog(a.config, "[AuthService] Login id: %s, platform: %s", id, platform)
	if id == "" {
		return nil, errors.New("用户ID不能为空")
	}
	strategies := map[string]func(string) (interface{}, error){
		PlatformWeb:         a.userClient.UserAuth,
		PlatformApp:         a.userClient.UserAuthToken,
		PlatformMiniProgram: a.userClient.UserAuthMiniProgramToken,
	}
	auth, ok := strategies[platform]
	if !ok {
		return nil, fmt.Errorf("不支持的平台类型: %s", platform)
	}
	info, err := auth(id)
	helper.DebugLog(a.config, "[AuthService] Login info: %+v, err: %v", info, err)
	if err != nil {
		return nil, fmt.Errorf("认证服务不可用: %w", err)
	}
	return info, nil
}

// LoginLogic 登录逻辑，生成token并存储用户信息
func (a *AuthService) LoginLogic(userInfo interface{}) (string, error) {
	helper.DebugLog(a.config, "[AuthService] LoginLogic userInfo: %+v", userInfo)
	res, ok := userInfo.(user_model.LoginRes)
	if !ok {
		helper.DebugLog(a.config, "[AuthService] LoginLogic 类型断言失败: %T", userInfo)
		return "", fmt.Errorf("无效的用户认证类型: %T", userInfo)
	}
	userData := res.LoginData
	if userData.ModuleToken == "" {
		userData.ModuleToken = helper.GenerateModuleToken(userData.ID, a.config.Skittle.Namespace)
		helper.DebugLog(a.config, "[AuthService] LoginLogic 生成新ModuleToken: %s", userData.ModuleToken)
	}
	if _, err := a.user.SetUserInfo(&userData); err != nil {
		helper.DebugLog(a.config, "[AuthService] LoginLogic 用户信息存储失败: %v", err)
		return "", fmt.Errorf("用户信息存储失败: %w", err)
	}
	helper.DebugLog(a.config, "[AuthService] LoginLogic 返回ModuleToken: %s", userData.ModuleToken)
	return userData.ModuleToken, nil
}
