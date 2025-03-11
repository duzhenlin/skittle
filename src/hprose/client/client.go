// Package hprose
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 18:39

package client

import (
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/core/helper"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/hprose/hprose-golang/rpc/fasthttp"
	"reflect"
)

// Client 单例
type Client struct {
	config     *config.Config
	ClientName string // 客户端名称
	isTo       bool
}

var client *Client

func init() {
	client = &Client{}
}

// GetClientInstance 获取实例
func GetClientInstance(config *config.Config) *Client {
	client.config = config
	return client
}

func (c *Client) GetConfig() *config.Config {
	return c.config
}

func (c *Client) SetConfig(config *config.Config) {
	c.config = config
}
func (c *Client) To(clientName string) *Client {
	c.ClientName = clientName
	c.isTo = true
	return c
}

func (c *Client) GetClient() (httpClient *fasthttp.FastHTTPClient, err error) {
	var UserServiceUrl string
	clients := c.config.Skittle.Client
	if c.ClientName == "" || c.isTo == false {
		c.ClientName = clients[0].ClientName
	}
	clientName := c.ClientName
	for _, item := range clients {
		if item.ClientName == clientName {
			UserServiceUrl = item.ClientUrl
			break
		}
	}

	httpClient = fasthttp.NewFastHTTPClient(UserServiceUrl)
	httpClient.AddInvokeHandler(c.clientAesInvokeHandler)
	return httpClient, err
}

func (c *Client) Invoke(name string, args []reflect.Value) (results []reflect.Value, err error) {
	httpClient, err := c.GetClient()
	if err != nil {
		return nil, err
	}
	settings := &rpc.InvokeSettings{
		Mode: rpc.Normal,
	}

	return httpClient.Invoke(name, args, settings)
}

func (c *Client) clientAesInvokeHandler(
	name string,
	args []reflect.Value,
	context rpc.Context,
	next rpc.NextInvokeHandler) (results []reflect.Value, err error) {
	newArgs := make([]reflect.Value, 2)

	newArgs[0] = reflect.ValueOf(helper.AesEncrypt(args[0], c.config.Skittle.SecretKey))
	newArgs[1] = reflect.ValueOf(c.config.Skittle.ModuleId)
	results, err = next(name, newArgs, context)
	return
}
