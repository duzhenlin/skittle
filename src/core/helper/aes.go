package helper

import (
	base64encoding "encoding/base64"
	"fmt"
	"github.com/cristalhq/base64"

	"github.com/duzhenlin/skittle/src/errors"
	"github.com/forgoer/openssl"
	"reflect"
	"strings"
)

func AesEncrypt(str reflect.Value, secretKey string) reflect.Value {
	origData := []byte(str.String()) // 待加密的数据
	key := []byte(secretKey)         // 加密的密钥
	dst, err := openssl.AesECBEncrypt(origData, key, openssl.PKCS7_PADDING)
	if err != nil {
		fmt.Printf("encrypt error:%s", err)
	}
	base64Str := base64.StdEncoding.EncodeToString(dst)
	encrypt := reflect.ValueOf(base64Str) // a reflect.Value
	return encrypt
}

func AesDecrypt(str string, secretKey string) (string, error) {
	//base64 解密处理原始字符串
	base64Dst := Base64DecryptToByte(str)
	if len(base64Dst) == 0 {
		return "", errors.ErrBase64DecryptToByte
	}
	//fmt.Printf("base64原始结果：%v \n", base64Dst)
	//fmt.Printf("base64结果：%s \n", string(base64Dst))
	key := []byte(secretKey) // 加密的密钥
	//aes ECB 16位 解密
	dst, err := openssl.AesECBDecrypt(base64Dst, key, openssl.PKCS7_PADDING)
	//fmt.Printf("aes结果：%s \n", dst)
	if err != nil {
		fmt.Println(err)
	}
	return string(dst), nil
}
func AesDecryptToHprose(str string, secretKey string) string {
	reader := strings.NewReader(str)
	decoder := base64encoding.NewDecoder(base64encoding.StdEncoding, reader)
	// base64以流式解码
	buf := make([]byte, 2)
	// base64保存解码后的数据
	base64Str := ""
	for {
		n, err := decoder.Read(buf)
		if n == 0 || err != nil {
			break
		}
		base64Str += string(buf[:n])
	}
	key := []byte(secretKey) // 加密的密钥
	dst, _ := openssl.AesECBDecrypt([]byte(base64Str), key, openssl.PKCS7_PADDING)
	return string(dst)
}
