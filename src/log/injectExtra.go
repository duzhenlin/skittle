// Package log
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/9
// @Time: 10:55

package log

import "context"

type ctxKeyExtra struct{}

func InjectExtraToContext(ctx context.Context, extra Extra) context.Context {
	return context.WithValue(ctx, ctxKeyExtra{}, extra)
}

func ExtractExtraFromContext(ctx context.Context) (Extra, bool) {
	val := ctx.Value(ctxKeyExtra{})
	extra, ok := val.(Extra)
	return extra, ok
}
