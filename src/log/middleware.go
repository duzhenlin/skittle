// Package log
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 15:49

package log

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func (l *DefaultLogger) ListenMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := l.LogRequest(r.Context(), r)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func (l *DefaultLogger) LogRequest(ctx context.Context, r *http.Request) context.Context {
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	ctx, extra := l.CoreLogger(ctx, r, bodyBytes)
	l.LogToES(ctx, INFO, "auto logger", extra)
	return ctx
}
