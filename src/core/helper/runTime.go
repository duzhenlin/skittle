package helper

import (
	"log"
	"time"
)

func RunTime(start time.Time, url string, name string, args map[string]interface{}) {
	elapsed := time.Since(start)
	log.Printf("该函数 （%v） %v ，参数 %v 执行完成耗时：%v\n", url, name, args, elapsed)
}
