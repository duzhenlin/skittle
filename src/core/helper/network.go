// Package helper
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 09:37

package helper

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

func GetClientIp(r *http.Request) string {
	var clientIp string
	realIps := r.Header.Get("X-Forwarded-For")
	if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
		ipArray := strings.Split(realIps, ",")
		clientIp = ipArray[0]
	}
	if clientIp == "" {
		clientIp = r.Header.Get("Proxy-Client-IP")
	}
	if clientIp == "" {
		clientIp = r.Header.Get("WL-Proxy-Client-IP")
	}
	if clientIp == "" {
		clientIp = r.Header.Get("HTTP_CLIENT_IP")
	}
	if clientIp == "" {
		clientIp = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if clientIp == "" {
		clientIp = r.Header.Get("X-Real-IP")
	}
	if clientIp == "" {
		clientIp = ExtractIPAddress(r.RemoteAddr)
	}
	return clientIp
}

// ExtractIPAddress 提取 IP 地址
func ExtractIPAddress(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// 如果解析失败，返回原始地址
		return remoteAddr
	}
	return host
}

// HeadersToJson 请求头转json
func HeadersToJson(header http.Header) string {
	b, _ := json.Marshal(header)
	return string(b)
}

// ExtractFormData 提取表单数据（支持multipart和普通表单）
func ExtractFormData(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err == nil && r.MultipartForm != nil {
			kv := map[string]interface{}{}
			for k, v := range r.MultipartForm.Value {
				if len(v) == 1 {
					kv[k] = v[0]
				} else {
					kv[k] = v
				}
			}
			for k, files := range r.MultipartForm.File {
				var names []string
				for _, f := range files {
					names = append(names, f.Filename)
				}
				if len(names) == 1 {
					kv[k] = names[0]
				} else if len(names) > 1 {
					kv[k] = names
				}
			}
			b, _ := json.Marshal(kv)
			return string(b)
		}
	} else if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		if err := r.ParseForm(); err == nil && len(r.PostForm) > 0 {
			b, _ := json.Marshal(r.PostForm)
			return string(b)
		}
	}
	return ""
}

// GetRequestPayload 获取请求体内容（formData优先，若有formData则body为空）
func GetRequestPayload(body []byte, formData string) string {
	if formData != "" {
		return ""
	}
	return string(body)
}
