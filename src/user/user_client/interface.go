// Package client
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 13:28

package user_client

type IUserClient interface {
	UserAuth(id string) (interface{}, error)
	UserAuthToken(token string) (interface{}, error)
	UserAuthMiniProgramToken(token string) (interface{}, error)
}
