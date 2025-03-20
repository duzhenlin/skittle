// Package hprose
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 18:39

package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/duzhenlin/skittle/src/core/helper/aes"
	"reflect"
	"sync"

	"github.com/duzhenlin/skittle/src/config"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/hprose/hprose-golang/rpc/fasthttp"
)

// HproseClient 定义客户端接口
type HproseClient interface {
	WithTarget(clientName string) *Client
	GetClient() (*fasthttp.FastHTTPClient, error)
	Invoke(method string, args ...interface{}) ([]reflect.Value, error)
}

// Client 表示Hprose客户端
type Client struct {
	config      *config.Config
	ctx         context.Context
	TargetName  string
	clientLock  sync.RWMutex
	initialized bool
}

var (
	instance *Client
	once     sync.Once
)

// NewClient 创建新的客户端实例（线程安全单例）
func NewClient(ctx context.Context, config *config.Config) *Client {
	once.Do(func() {
		instance = &Client{
			config:      config,
			ctx:         ctx,
			initialized: true,
		}
	})
	return instance
}

// WithTarget 设置目标客户端（返回新实例保证线程安全）
func (c *Client) WithTarget(clientName string) *Client {
	c.clientLock.Lock()
	defer c.clientLock.Unlock()
	c.TargetName = clientName
	return &Client{
		config:      c.config,
		ctx:         c.ctx,
		TargetName:  clientName,
		initialized: true,
	}
}

// getServiceURL 获取服务地址（带错误处理）
func (c *Client) getServiceURL() (string, error) {
	if c.config == nil || c.config.Skittle.Client == nil {
		return "", errors.New("client configuration not initialized")
	}

	clients := *c.config.Skittle.Client
	if len(clients) == 0 {
		return "", errors.New("empty client configuration")
	}

	target := c.TargetName
	if target == "" {
		target = clients[0].ClientName
	}

	for _, client := range clients {
		if client.ClientName == target {
			return client.ClientUrl, nil
		}
	}

	return "", fmt.Errorf("target client not found: %s", target)
}

// GetClient 获取Hprose客户端实例
func (c *Client) GetClient() (*fasthttp.FastHTTPClient, error) {
	if !c.initialized {
		return nil, errors.New("client not initialized")
	}

	serviceURL, err := c.getServiceURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get service URL: %w", err)
	}

	client := fasthttp.NewFastHTTPClient(serviceURL)
	client.AddInvokeHandler(c.clientAesInvokeHandler)
	return client, nil
}

// clientAesInvokeHandler AES加密中间件
func (c *Client) clientAesInvokeHandler(
	name string,
	args []reflect.Value,
	ctx rpc.Context,
	next rpc.NextInvokeHandler) ([]reflect.Value, error) {

	if len(args) == 0 {
		return nil, errors.New("encryption requires at least one argument")
	}
	encryptedArgs := make([]reflect.Value, 2)
	encrypt, err := aes.Encrypt(args[0].String(), c.config.Skittle.SecretKey)
	if err != nil {
		return nil, err
	}
	encryptedArgs[0] = reflect.ValueOf(encrypt)
	encryptedArgs[1] = reflect.ValueOf(c.config.Skittle.ModuleId)

	return next(name, encryptedArgs, ctx)
}

// Invoke 执行远程调用
func (c *Client) Invoke(method string, args ...interface{}) ([]reflect.Value, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, fmt.Errorf("client initialization failed: %w", err)
	}

	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}

	return client.Invoke(method, reflectArgs, &rpc.InvokeSettings{Mode: rpc.Normal})
}
