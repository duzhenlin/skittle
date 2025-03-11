package helper

import (
	"fmt"
	"github.com/tjfoc/gmsm/sm3"
)

func Sm3(str string) string {
	h := sm3.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	//将[]byte转成16进制
	return fmt.Sprintf("%x", sum)
}
