// Package hprose
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 18:39

package hprose

import (
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/core/helper"
	"github.com/hprose/hprose-golang/rpc"
	"reflect"
)

type Server struct {
	config *config.Config
}

var server *Server

func init() {
	server = &Server{}
}

// GetServerInstance 获取实例
func GetServerInstance(config *config.Config) *Server {
	server.config = config
	return server
}

func (s *Server) GetConfig() *config.Config {
	return s.config
}

func (s *Server) Start() (service *rpc.HTTPService) {
	service = rpc.NewHTTPService()
	service.AddInvokeHandler(s.serverAesInvokeHandler)

	if s.config.Skittle.Server.IsModule {
		service.AddFunction("notice", func() {})
		service.AddFunction("register", func() {})
	}

	return service
}

// serverAesInvokeHandler 服务端中间件，主要复制解密
func (s *Server) serverAesInvokeHandler(
	name string,
	args []reflect.Value,
	context rpc.Context,
	next rpc.NextInvokeHandler,
) (results []reflect.Value, err error) {
	//转成字符串
	toString := args[0].Interface().(string)
	newArgs := make([]reflect.Value, 2)
	//aes解密
	decryptStr := helper.AesDecryptToHprose(toString, s.GetConfig().Skittle.SecretKey)
	newArgs[0] = reflect.ValueOf(decryptStr)
	newArgs[1] = reflect.ValueOf("module_id")
	results, err = next(name, newArgs, context)
	return
}

func CallMethodByName(obj interface{}, methodName string) {
	v := reflect.ValueOf(obj)
	method := v.MethodByName(methodName)
	if method.IsValid() {
		method.Call(nil)
	} else {
		fmt.Printf("Method %s not found.\n", methodName)
	}
}

func CallMethodByNameWithArgs(obj interface{}, methodName string, args ...interface{}) []interface{} {
	// 获取对象的反射值
	v := reflect.ValueOf(obj)
	// 通过方法名获取方法的反射值
	method := v.MethodByName(methodName)

	if method.IsValid() {
		// 检查传入的参数数量是否与方法所需的参数数量匹配
		if len(args) != method.Type().NumIn() {
			fmt.Printf("Incorrect number of arguments for method %s. Expected %d, got %d.\n",
				methodName, method.Type().NumIn(), len(args))
			return nil
		}

		// 将传入的参数转换为 reflect.Value 类型的切片
		var reflectArgs []reflect.Value
		for _, arg := range args {
			reflectArgs = append(reflectArgs, reflect.ValueOf(arg))
		}

		// 调用方法并获取返回值
		results := method.Call(reflectArgs)

		// 将反射返回值转换为 interface{} 类型的切片
		var resultValues []interface{}
		for _, result := range results {
			resultValues = append(resultValues, result.Interface())
		}
		return resultValues
	}
	fmt.Printf("Method %s not found.\n", methodName)
	return nil
}
