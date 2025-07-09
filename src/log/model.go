// Package log_client
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 09:38

package log

import "time"

type User struct {
	Id           string `json:"id"`
	Nickname     string `json:"nickname"`
	Organization string `json:"organization"`
	Phone        string `json:"phone"`
}
type General struct {
	RequestUrl    string `json:"request_url"`
	RequestMethod string `json:"request_method"`
	RemoteAddress string `json:"remote_address"`
}
type Http struct {
	General        General `json:"general"`
	RequestHeaders string  `json:"request_headers"`
	RequestPayload string  `json:"request_payload"`
	FormData       string  `json:"form_data"`
}

type Extra struct {
	RuntimeUuid string `json:"runtime_uuid"`
	PHP         Php    `json:"php"`
	User        User   `json:"user"`
	Http        Http   `json:"http"`
}
type Php struct {
	Version int    `json:"version"`
	Sapi    string `json:"sapi"`
}

type LogEntry struct {
	Message   string    `json:"message"`
	Context   string    `json:"context"`
	Level     int       `json:"level"`
	LevelName string    `json:"level_name"`
	Channel   string    `json:"channel"`
	Datetime  time.Time `json:"datetime"`
	Extra     Extra     `json:"extra"`
}
