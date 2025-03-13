// Package config
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 16:40

package config

type Config struct {
	// 配置字段
	Redis   Redis
	Skittle Skittle
	Debug   bool `default:"false"`
}
type Skittle struct {
	Namespace string   `json:"namespace"`
	Server    Server   `json:"server"`
	Client    []Client `json:"client"`
	ModuleId  string   `json:"module_id"`
	SecretKey string   `json:"secret_key"`
}

type Client struct {
	ClientName string `json:"client_name"`
	ClientType string `json:"client_type"`
	ClientUrl  string `json:"client_url"`
	ClientTag  string `json:"client_tag"`
}

type Server struct {
	IsModule bool        `default:"false"`
	Register interface{} `default:""`
}

type Redis struct {
	Cont string `default:"127.0.0.1:6379"`
	Pwd  string `default:""`
	Db   int    `default:"0"`
}

func (c *Config) SetConfig(config Config) {
	*c = config
}
