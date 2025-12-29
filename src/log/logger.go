// Package log
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 15:56

package log

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/duzhenlin/skittle/v2/src/user/user_service"

	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/duzhenlin/skittle/v2/src/core/helper"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...interface{})
	Info(ctx context.Context, msg string, fields ...interface{})
	Warn(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, fields ...interface{})
	LogToES(ctx context.Context, level int, msg string, extra interface{})
	ListenMiddleware() func(http.Handler) http.Handler
	LogRequest(ctx context.Context, r *http.Request) context.Context
	CoreLogger(ctx context.Context, r *http.Request, body []byte) (context.Context, Extra)
}

type DefaultLogger struct {
	Ctx       context.Context
	cfg       *config.Config
	zapLogger *ZapLogger
	esLogger  *ESLogger
	User      *user_service.UserService
}

func NewDefaultLogger(deps LogDeps, zapLogger *ZapLogger, esLogger *ESLogger) *DefaultLogger {
	return &DefaultLogger{
		zapLogger: zapLogger,
		esLogger:  esLogger,
		cfg:       deps.Config,
		Ctx:       deps.Ctx,
		User:      deps.User,
	}
}

func (l *DefaultLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.zapLogger.Debug(msg, fields...)
	extra, _ := ExtractExtraFromContext(ctx)
	l.LogToES(ctx, DEBUG, msg, extra)
}
func (l *DefaultLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.zapLogger.Info(msg, fields...)
	extra, _ := ExtractExtraFromContext(ctx)
	l.LogToES(ctx, INFO, msg, extra)
}
func (l *DefaultLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.zapLogger.Warn(msg, fields...)
	extra, _ := ExtractExtraFromContext(ctx)
	l.LogToES(ctx, WARNING, msg, extra)
}
func (l *DefaultLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.zapLogger.Error(msg, fields...)
	extra, _ := ExtractExtraFromContext(ctx)
	l.LogToES(ctx, ERROR, msg, extra)
}
func (l *DefaultLogger) LogToES(ctx context.Context, level int, msg string, extra interface{}) {
	channel := ""
	if l.cfg != nil && l.cfg.Skittle != nil {
		channel = l.cfg.Skittle.Namespace
	}
	doc := map[string]interface{}{
		"message":    msg,
		"level":      level,
		"level_name": LevelName(level),
		"channel":    channel,
		"datetime":   time.Now(),
		"extra":      mergeExtraWithContext(ctx, extra),
	}
	if l.esLogger != nil {
		err := l.esLogger.WriteLog(doc)
		if err != nil {
			log.Printf("write log to es error: %v", err)
			return
		}
	}
}

func (l *DefaultLogger) CoreLogger(ctx context.Context, r *http.Request, body []byte) (context.Context, Extra) {
	// 1. 获取用户信息（如果user模块可用）
	var user User
	if l.User != nil {
		userData, err := l.User.GetUserInfoNew(r)
		if err == nil {
			user = User{
				Id:           userData.ID,
				Nickname:     userData.Nickname,
				Organization: userData.Orgname,
				Phone:        userData.Username,
			}
		}
	}

	// 2. 通用辅助方法
	formData := helper.ExtractFormData(r)
	httpInfo := Http{
		General: General{
			RequestUrl:    r.URL.Path,
			RequestMethod: r.Method,
			RemoteAddress: helper.GetClientIp(r),
		},
		RequestHeaders: helper.HeadersToJson(r.Header),
		RequestPayload: helper.GetRequestPayload(body, formData),
		FormData:       formData,
	}

	// 3. 通过 helper 方法自动提取 runtimeUuid
	runtimeUuid := helper.GetRuntimeUUIDFromContext(ctx)

	extra := Extra{
		RuntimeUuid: runtimeUuid,
		PHP:         Php{Version: helper.GetGoVersionInt(), Sapi: "golang"},
		User:        user,
		Http:        httpInfo,
	}
	ctx = InjectExtraToContext(ctx, extra) // 注入 extra
	return ctx, extra
}

// mergeExtraWithContext 合并extra和context中的请求信息
func mergeExtraWithContext(ctx context.Context, extra interface{}) interface{} {
	// 这里直接返回 extra，后续如需合并 context 可扩展
	return extra
}
