// Package helper
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/4
// @Time: 15:27

package helper

import (
	"github.com/hprose/hprose-golang/util"
	"strings"
)

// GenerateModuleToken 生成模块访问令牌
func GenerateModuleToken(userID, namespace string) string {
	builder := strings.Builder{}
	builder.WriteString(util.UUIDv4())
	builder.WriteString(userID)
	builder.WriteString(namespace)
	return GetStringMd5(builder.String())
}
