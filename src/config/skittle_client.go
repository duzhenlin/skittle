// Package config
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/20
// @Time: 01:16

package config

type Client struct {
	ClientName string `json:"client_name"`
	ClientType string `json:"client_type"`
	ClientUrl  string `json:"client_url"`
	ClientTag  string `json:"client_tag"`
}
