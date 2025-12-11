// Package core
// 版本信息管理
// 作为依赖包，只提供包的版本号和运行时信息
// 构建信息（构建日期、Git提交等）应由主程序管理
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com

package core

import (
	"fmt"
	"runtime"
)

const (
	// Version 当前版本号
	// 遵循语义化版本规范: https://semver.org/
	Version = "2.0.4"
)

// VersionInfo 版本信息结构
// 只包含依赖包本身的版本信息和运行时环境信息
type VersionInfo struct {
	Version   string `json:"version"`    // 依赖包版本号
	GoVersion string `json:"go_version"` // Go运行时版本
	Platform  string `json:"platform"`   // 操作系统平台
	Arch      string `json:"arch"`       // CPU架构
}

// GetVersion 获取版本号
func GetVersion() string {
	return Version
}

// GetVersionInfo 获取版本信息
// 返回依赖包的版本信息和运行时环境信息
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   Version,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// GetVersionString 获取格式化的版本信息字符串
func GetVersionString() string {
	info := GetVersionInfo()
	return fmt.Sprintf(
		"Skittle v%s\n"+
			"  Go Version: %s\n"+
			"  Platform: %s/%s",
		info.Version,
		info.GoVersion,
		info.Platform,
		info.Arch,
	)
}

// GetShortVersion 获取简短版本信息
func GetShortVersion() string {
	return fmt.Sprintf("Skittle v%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)
}
