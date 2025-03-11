package errors

import (
	"errors"
)

var (
	ErrBase64DecryptToByte = errors.New("base64 解码失败")
)
