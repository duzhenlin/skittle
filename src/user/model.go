// Package user
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 22:13

package user

import jsoniter "github.com/json-iterator/go"

type LoginData struct {
	Orgname  string `json:"orgname"`
	Isleader string `json:"isleader"`
	Content  struct {
		ColumnRead []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Key  string `json:"key"`
			Pid  int    `json:"pid"`
		} `json:"column_read"`
		ColumnWrite []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Key  string `json:"key"`
			Pid  int    `json:"pid"`
		} `json:"column_write"`
	} `json:"content"`
	ID              string `json:"id"`
	Loginname       string `json:"loginname"`
	Nickname        string `json:"nickname"`
	Orgid           string `json:"orgid"`
	IsCorrespondent string `json:"is_correspondent"`
	Username        string `json:"username"`
	Roleid          string `json:"roleid"`
	Priv            any    `json:"priv"`
	Department      string `json:"department"`
	Orgopenid       string `json:"orgopenid"`
	Avatar          string `json:"avatar"`
	ModuleToken     string `json:"module_token"`
	Dptname         string `json:"dptname"`
}

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
