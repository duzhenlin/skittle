package cache

import (
	"github.com/duzhenlin/skittle/v2/src/config"
)

// NewCache 工厂方法，根据配置返回 Cache 实例
func NewCache(cfg *config.Config) (Cache, error) {
	// 目前仅支持 redis，可扩展 memory、其他实现
	return NewRedisCache(cfg)
}
