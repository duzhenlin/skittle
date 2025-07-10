// Package redis
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/10
// @Time: 22:27

package cache

// Cache 缓存接口，便于扩展多种实现
// value 建议为 string，复杂结构建议自行序列化
// ttl 单位为秒，0 表示不过期

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl int) error
	Del(key string) error
	Exists(key string) (bool, error)
	Expire(key string, ttl int) error // 新增
}
