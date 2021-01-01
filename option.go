package httpoet

import (
	"context"
	"time"
)

type Option func(req *RequestBuilder) (fnDefer func())

func OBackground() Option {
	return func(*RequestBuilder) func() { return func() {} }
}

func OTimeout(d time.Duration) Option {
	if d <= 0 {
		return OBackground()
	}
	return func(req *RequestBuilder) func() {
		orgCtx := req.Context
		if orgCtx == nil {
			orgCtx = context.Background()
		}
		ctxWithTimeout, cancel := context.WithTimeout(orgCtx, d)
		req = req.XContext(ctxWithTimeout)
		return cancel
	}
}

func OAddHeaders(headers H) Option {
	return func(req *RequestBuilder) func() {
		req.Header.Append(headers)
		return func() {}
	}
}

func OSetHeader(key string, value ...string) Option {
	return func(req *RequestBuilder) func() {
		req.Header.Set(key, value...)
		return func() {}
	}
}

func OAppendQuery(key string, val string) Option {
	return func(req *RequestBuilder) func() {
		req.Query.Add(key, val)
		return func() {}
	}
}

func OAppendQueryH(queries Q) Option {
	return func(req *RequestBuilder) func() {
		req.Query.Append(queries)
		return func() {}
	}
}
