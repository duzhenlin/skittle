// Package config
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/20
// @Time: 01:14

package config

type Redis struct {
	ConType      string   `default:"tcp"`
	TpcConArr    string   `default:"127.0.0.1:6379"`
	Pwd          string   `default:""`
	Db           int      `default:"0"`
	MasterName   string   `default:""`
	SentinelAddr []string `default:""`
	PoolSize     int      `default:"50"` // 设置最大连接数
}
