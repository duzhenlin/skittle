// Package config
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:40

package config

type Config struct {
	// 配置字段
	Debug    bool `default:"false"`
	Redis    *Redis
	Skittle  *Skittle
	EsConfig *EsConfig
}
type Skittle struct {
	Namespace string    `json:"namespace"`
	Server    *Server   `json:"server"`
	Client    *[]Client `json:"client"`
	ModuleId  string    `json:"module_id"`
	SecretKey string    `json:"secret_key"`
}

type EsConfig struct {
	Address   string `json:"address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	IndexName string `json:"index_name"`
	Channel   string `json:"channel"`
}

func (c *Config) SetConfig(config Config) {
	*c = config
}
