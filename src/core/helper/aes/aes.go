// Package aes
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/19
// @Time: 22:39

package aes

import (
	"fmt"
	"github.com/cristalhq/base64"
	"github.com/duzhenlin/skittle/v2/src/constant"
	"github.com/duzhenlin/skittle/v2/src/core/helper"
	"github.com/duzhenlin/skittle/v2/src/errors"
	"github.com/forgoer/openssl"
)

const (
	// AES标准密钥长度

	KeyLen128 = 16
	KeyLen192 = 24
	KeyLen256 = 32
)

// Encrypt 使用AES-ECB模式加密数据
// plaintext: 明文字符串
// secretKey: 加密密钥（16/24/32字节）
func Encrypt(plaintext, secretKey string) (string, error) {
	if err := validateKey(secretKey); err != nil {
		return "", fmt.Errorf("encryption validation failed: %w", err)
	}

	ciphertext, err := openssl.AesECBEncrypt([]byte(plaintext), []byte(secretKey), openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 使用AES-ECB模式解密数据
// ciphertext: Base64编码的密文字符串
// secretKey: 加密密钥（16/24/32字节）
func Decrypt(ciphertext, secretKey string, base64Type string) (string, error) {
	if err := validateKey(secretKey); err != nil {
		return "", fmt.Errorf("decryption validation failed: %w", err)
	}
	// 统一解码处理
	var decoded []byte
	var decodeErr error

	switch base64Type {
	case constant.Base64Standard:
		decoded, decodeErr = helper.Base64Decode(ciphertext)
	case constant.Base64URL:
		decoded, decodeErr = helper.Base64URLDecode(ciphertext)
	default:
		return "", fmt.Errorf("不支持的Base64类型: %s", base64Type)
	}

	if decodeErr != nil {
		return "", fmt.Errorf("base64 decode failed: %w", decodeErr)
	}
	if len(decoded) == 0 {
		return "", errors.ErrBase64DecodeFailed
	}

	plaintext, err := openssl.AesECBDecrypt(decoded, []byte(secretKey), openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errors.ErrDecryptionFailed, err)
	}

	return string(plaintext), nil
}

// 密钥校验函数
func validateKey(key string) error {
	switch len(key) {
	case KeyLen128, KeyLen192, KeyLen256:
		return nil
	default:
		return errors.ErrInvalidKeyLength
	}
}
