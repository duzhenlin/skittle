// Package errors
// 统一错误定义
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8

package errors

import (
	"fmt"
)

// 应用级错误
var (
	// ErrConfigInvalid 配置无效
	ErrConfigInvalid = New("config invalid")

	// ErrModuleNotEnabled 模块未启用
	ErrModuleNotEnabled = New("module not enabled")

	// ErrServiceUnavailable 服务不可用
	ErrServiceUnavailable = New("service unavailable")

	// ErrDependencyMissing 依赖缺失
	ErrDependencyMissing = New("dependency missing")
)

// AppError 应用错误结构
type AppError struct {
	Code    string // 错误代码
	Message string // 错误消息
	Cause   error  // 原始错误
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 返回原始错误（支持 errors.Unwrap）
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New 创建新的应用错误
func New(message string) *AppError {
	return &AppError{
		Code:    "SKITTLE_ERROR",
		Message: message,
	}
}

// NewWithCode 创建带错误代码的应用错误
func NewWithCode(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(err error, message string) *AppError {
	return &AppError{
		Code:    "SKITTLE_ERROR",
		Message: message,
		Cause:   err,
	}
}

// WrapWithCode 使用错误代码包装错误
func WrapWithCode(err error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// Is 检查错误是否匹配
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}

