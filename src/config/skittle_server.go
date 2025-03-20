// Package config
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/20
// @Time: 01:14

package config

type NoticeFunctionFunc func(string, string, string)

type Server struct {
	IsModule       bool      `default:"false"`
	Register       *Register `default:""`
	FunctionList   []string
	UserUpdateFunc interface{} `yaml:"user_update_func"`
}

type Register struct {
	Data    *[]RegisterData    `json:"data"`
	Content *[]RegisterContent `json:"content"`
	Quick   *[]RegisterQuick   `json:"quick"`
}
type RegisterOps struct {
	URL    string `json:"url"`
	Method int    `json:"method"`
}

type RegisterData struct {
	ID   string        `json:"id"`
	Pid  int           `json:"pid"`
	Name string        `json:"name"`
	Ops  []RegisterOps `json:"ops"`
}

type RegisterContent struct {
	Type    string `json:"type"`
	Key     string `json:"key"`
	KeyName string `json:"key_name"`
}
type RegisterQuick struct {
	Name string `json:"name"`
	Icon int    `json:"icon"`
	Path string `json:"path"`
}
