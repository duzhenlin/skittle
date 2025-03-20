// Package config
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:40

package config

type Config struct {
	// 配置字段
	Debug   bool `default:"false"`
	Redis   *Redis
	Skittle *Skittle
}
type Skittle struct {
	Namespace string    `json:"namespace"`
	Server    *Server   `json:"server"`
	Client    *[]Client `json:"client"`
	ModuleId  string    `json:"module_id"`
	SecretKey string    `json:"secret_key"`
}

func (c *Config) SetConfig(config Config) {
	*c = config
}
