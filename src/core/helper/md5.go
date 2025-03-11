package helper

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
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

func MakeSign(id string, random string, signTime string, secret string) string {
	//签名时间戳
	signTime64, _ := strconv.ParseInt(signTime, 10, 64)
	//当前时间戳
	NowTime64 := time.Now().Unix()
	//当前时间往后延长五秒有效期,并且签名时间需要比当前时间相等或者大
	if (signTime64+5 < NowTime64) && (signTime64 < NowTime64) {
		return "失败"
	} else {
		str := id + random + secret + signTime
		return GetStringMd5(str)
	}

}
