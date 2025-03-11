// Package client
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/11
// @Time: 12:02

package client

type ServiceDefaultUser struct {
	UserAuth      func(string, string) (interface{}, error) `name:"user_auth"`
	UserAuthToken func(string, string) (interface{}, error) `name:"user_auth_token"`
}

func (c *Client) GetDefaultUserService() ServiceDefaultUser {
	var serviceDefaultUser ServiceDefaultUser
	c.SetConfig(c.config)
	c.To("user")
	serviceClient, err := c.GetClient()
	if err != nil {
		panic(err)
	}
	serviceClient.UseService(&serviceDefaultUser)
	return serviceDefaultUser
}
