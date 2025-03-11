package helper

import (
	"fmt"
	"github.com/cristalhq/base64"
	"github.com/duzhenlin/skittle/src/core/helper/base64url"
)

// Base64DecryptToByte   base解密 /*
func Base64DecryptToByte(str string) []byte {
	base64Str, err := base64url.Decode(str)
	if err != nil {
		fmt.Println("error:", err)
	}
	maxLen := base64.StdEncoding.DecodedLen(len(base64Str))
	dst := make([]byte, maxLen)
	n, err := base64.StdEncoding.Decode(dst, base64Str)
	if err != nil {
		fmt.Println("error1:", err)
	}
	return dst[:n]
}
