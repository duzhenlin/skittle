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

func (c *Client) GetDefaultUserService() (ServiceDefaultUser, error) {
	var serviceDefaultUser ServiceDefaultUser
	c.WithTarget("user")
	serviceClient, err := c.GetClient()
	if err != nil {
		return serviceDefaultUser, err
	}
	serviceClient.UseService(&serviceDefaultUser)
	return serviceDefaultUser, nil
}
