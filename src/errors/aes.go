// Package errors
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/3/19
// @Time: 22:06

package errors

import "errors"

// AES加密方法统一错误处理
var (
	ErrBase64DecodeFailed = errors.New("base64 decode failed")
	ErrDecryptionFailed   = errors.New("decryption failed")
	ErrInvalidKeyLength   = errors.New("invalid key length")
)
