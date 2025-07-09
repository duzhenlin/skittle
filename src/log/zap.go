// Package log
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 15:55

package log

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	l, _ := zap.NewProduction()
	return &ZapLogger{logger: l.Sugar()}
}

func (z *ZapLogger) Debug(msg string, fields ...interface{}) { z.logger.Debugw(msg, fields...) }
func (z *ZapLogger) Info(msg string, fields ...interface{})  { z.logger.Infow(msg, fields...) }
func (z *ZapLogger) Warn(msg string, fields ...interface{})  { z.logger.Warnw(msg, fields...) }
func (z *ZapLogger) Error(msg string, fields ...interface{}) { z.logger.Errorw(msg, fields...) }
