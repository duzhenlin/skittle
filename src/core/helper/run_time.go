package helper

import (
	"log"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

func RunTime(start time.Time, url string, name string, args map[string]interface{}) {
	elapsed := time.Since(start)
	log.Printf("该函数 （%v） %v ，参数 %v 执行完成耗时：%v\n", url, name, args, elapsed)
}

func GetGoVersionMajorMinor() (int, int) {
	re := regexp.MustCompile(`go(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(runtime.Version())
	if len(matches) >= 3 {
		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		return major, minor
	}
	return 0, 0
}

func GetGoVersionInt() int {
	re := regexp.MustCompile(`go(\d+)\.(\d+)(?:\.(\d+))?`)
	matches := re.FindStringSubmatch(runtime.Version())
	if len(matches) >= 3 {
		versionStr := matches[1] + matches[2]
		if len(matches) >= 4 && matches[3] != "" {
			versionStr += matches[3]
		}
		versionInt, _ := strconv.Atoi(versionStr)
		return versionInt
	}
	return 0
}
