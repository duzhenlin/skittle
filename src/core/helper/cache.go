package helper

import (
	"github.com/muesli/cache2go"
	"time"
)

// Add
//
//	@Description: 加入缓存
//	@param key
//	@param data
//	@param lifeSpan
//	@return *cache2go.CacheItem
func Add(key interface{}, data interface{}, lifeSpan time.Duration) *cache2go.CacheItem {
	cache := cache2go.Cache("myCache")
	return cache.Add(key, lifeSpan, data)
}

// Get
//
//	@Description: 获取缓存
//	@param key
//	@return *cache2go.CacheItem
//	@return error
func Get(key interface{}) (*cache2go.CacheItem, error) {
	cache := cache2go.Cache("myCache")
	return cache.Value(key)
}
