package helper

import (
	base64encoding "encoding/base64"
	"fmt"
	"github.com/cristalhq/base64"
	"github.com/duzhenlin/skittle/v2/src/core/helper/base64url"
	"io"
	"strings"
)

// Base64URLDecode Base64解码
// @param str Base64 URL编码字符串
// @return []byte 字节切片
// @return error 错误
func Base64URLDecode(str string) ([]byte, error) {

	base64Str, err := base64url.Decode(str)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}
	maxLen := base64.StdEncoding.DecodedLen(len(base64Str))
	dst := make([]byte, maxLen)
	n, err := base64.StdEncoding.Decode(dst, base64Str)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}
	return dst[:n], nil
}

// Base64Decode 解码Base64编码字符串
// 参数：
//
//	encodedStr - 标准Base64编码字符串（支持带填充）
//
// 返回值：
//
//	string - 解码后的原始字符串
//	error - 解码过程中的错误信息
func Base64Decode(encodedStr string) ([]byte, error) {
	// 使用标准解码器（自动处理填充）
	decoder := base64encoding.NewDecoder(base64encoding.StdEncoding, strings.NewReader(encodedStr))

	// 使用更高效的缓冲区
	buf := make([]byte, 4096) // 4KB缓冲区
	var result strings.Builder

	for {
		n, err := decoder.Read(buf)
		if n > 0 {
			result.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("base64解码失败: %w", err)
		}
	}
	return []byte(result.String()), nil
}
