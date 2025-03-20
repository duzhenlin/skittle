package helper

import (
	"crypto/md5"
	"encoding/hex"
)

// GetStringMd5
// @Description: 计算字符串md5
// @param s
// @return string
func GetStringMd5(s string) string {
	str := md5.New()
	str.Write([]byte(s))
	md5Str := hex.EncodeToString(str.Sum(nil))
	return md5Str
}
