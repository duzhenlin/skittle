// Package server
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/19
// @Time: 23:28

package hprose_server

import (
	"encoding/json"
	"github.com/duzhenlin/skittle/v2/src/config"
	"log"
)

type noticeParam struct {
	Id     string `json:"id"`
	Action string `json:"action"`
	OrgId  string `json:"orgid"`
}

// noticeFunction 用户更新函数
func (s *HproseServerService) noticeFunction(data string) *noticeParam {
	var param noticeParam
	if err := json.Unmarshal([]byte(data), &param); err != nil {
		log.Printf("JSON解析失败: %v", err)
		return nil
	}
	// 参数完整性检查
	if param.Id == "" || param.OrgId == "" || param.Action == "" {
		if s.config.Debug {
			log.Printf("参数缺失: id=%s, org_id=%s, action=%s", param.Id, param.OrgId, param.Action)
		}
		return &param
	}
	// 配置存在性检查
	if s.config.Skittle.Server.UserUpdateFunc == nil {
		if s.config.Debug {
			log.Printf("[DEBUG] user update function not configured")
		}
		return &param
	}

	// 双重类型检查
	switch fn := s.config.Skittle.Server.UserUpdateFunc.(type) {
	case config.NoticeFunctionFunc:
		fn(param.Id, param.OrgId, param.Action)
		if s.config.Debug {
			log.Printf("[DEBUG] 用户更新函数执行成功，参数: %+v", param)
		}
	case func(string, string, string):
		fn(param.Id, param.OrgId, param.Action)
		if s.config.Debug {
			log.Printf("[DEBUG] 用户更新函数执行成功，参数: %+v", param)
		}
	default:
		if s.config.Debug {
			log.Printf("[DEBUG] 用户更新函数类型错误，期望 func(string,string,string)，实际类型 %T",
				s.config.Skittle.Server.UserUpdateFunc)
		}
	}

	return &param
}

// registerFunction 注册函数
func (s *HproseServerService) registerFunction() *config.Register {
	// 防御性检查配置层级
	if s.config == nil ||
		s.config.Skittle == nil ||
		s.config.Skittle.Server == nil {
		return nil
	}

	// 配置存在性检查
	registerConfig := s.config.Skittle.Server.Register
	if registerConfig == nil {
		return nil
	}

	return &config.Register{
		Data:    registerConfig.Data,
		Quick:   registerConfig.Quick,
		Content: registerConfig.Content,
	}
}
