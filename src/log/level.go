// Package log
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 15:48

package log

const (
	DEBUG     = 100
	INFO      = 200
	WARNING   = 300
	ERROR     = 400
	CRITICAL  = 500
	ALERT     = 550
	EMERGENCY = 600
)

func LevelName(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	case ALERT:
		return "ALERT"
	case EMERGENCY:
		return "EMERGENCY"
	default:
		return "UNKNOWN"
	}
}
