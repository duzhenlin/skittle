// Package hprose
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 18:39

package server

import (
	"context"
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/constant"
	"github.com/duzhenlin/skittle/src/core/helper/aes"
	"github.com/hprose/hprose-golang/rpc"
	"reflect"
)

// HproseServer 定义客户端接口
type HproseServer interface {
	Start() *rpc.HTTPService
}

type Server struct {
	config *config.Config
	ctx    context.Context
}

// NewServer 创建新的服务实例
func NewServer(ctx context.Context, config *config.Config) *Server {
	return &Server{
		config: config,
		ctx:    ctx,
	}
}

// Start 启动服务
func (s *Server) Start() *rpc.HTTPService {
	service := rpc.NewHTTPService()
	service.AddInvokeHandler(s.serverAesInvokeHandler)
	// 注册服务
	if s.config.Skittle.Server.IsModule {
		// 注册默认服务
		service.AddFunction("notice", s.noticeFunction)
		service.AddFunction("register", s.registerFunction)
	}

	return service
}

// serverAesInvokeHandler 处理aes加密的请求
func (s *Server) serverAesInvokeHandler(
	name string,
	args []reflect.Value,
	context rpc.Context,
	next rpc.NextInvokeHandler,
) (results []reflect.Value, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("invalid arguments")
	}
	//转成字符串
	strArg, ok := args[0].Interface().(string)
	if !ok {
		return nil, fmt.Errorf("expected string argument, got %T", args[0].Interface())
	}
	decrypted, err := aes.Decrypt(strArg, s.config.Skittle.SecretKey, constant.Base64Standard)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	newArgs := []reflect.Value{
		reflect.ValueOf(decrypted),
		reflect.ValueOf(args[1]),
	}

	return next(name, newArgs, context)
}
