// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 22:13

package user

import jsoniter "github.com/json-iterator/go"

// ColumnPermission 定义通用的权限列结构
type ColumnPermission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	Pid  int    `json:"pid"`
}

// Content 定义用户内容权限
type Content struct {
	ColumnRead  []ColumnPermission `json:"column_read"`
	ColumnWrite []ColumnPermission `json:"column_write"`
}

// LoginData 用户登录数据
type LoginData struct {
	Orgname         string  `json:"orgname"`
	Isleader        string  `json:"isleader"`
	Content         Content `json:"content"`
	ID              string  `json:"id"`
	Loginname       string  `json:"loginname"`
	Nickname        string  `json:"nickname"`
	Orgid           string  `json:"orgid"`
	IsCorrespondent string  `json:"is_correspondent"`
	Username        string  `json:"username"`
	Roleid          string  `json:"roleid"`
	Priv            any     `json:"priv,omitempty"` // 使用omitempty忽略空值
	Department      string  `json:"department"`
	Orgopenid       string  `json:"orgopenid"`
	Avatar          string  `json:"avatar"`
	ModuleToken     string  `json:"module_token"`
	Dptname         string  `json:"dptname"`
}

// LoginSign 登录签名信息
type LoginSign struct {
	Token string `json:"token"`
}

type LoginRes struct {
	Code      int       `json:"code"`
	LoginData LoginData `json:"data"`
}

func (m *LoginData) MarshalBinary() (data []byte, err error) {
	return jsoniter.Marshal(m)
}

func (m LoginData) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, &m)
}
