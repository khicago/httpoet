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

func OSetHeaders(headers IHeader) Option {
	return func(req *RequestBuilder) func() {
		if req.Header == nil {
			req.Header = make(Hs)
		}

		req.Header = req.Header.WithH(headers)
		return func() {}
	}
}

func OAddHeaders(headers IHeader) Option {
	return func(req *RequestBuilder) func() {
		if req.Header == nil {
			req.Header = make(Hs)
		}
		req.Header = req.Header.WithHAppend(headers)
		return func() {}
	}
}

func OSetHeader(key string, value ...string) Option {
	return func(req *RequestBuilder) func() {
		if req.Header == nil {
			req.Header = make(Hs)
		}

		req.Header = req.Header.WithKV(key, value...)
		return func() {}
	}
}

func OAppendHeader(key string, value ...string) Option {
	return func(req *RequestBuilder) func() {
		if req.Header == nil {
			req.Header = make(Hs)
		}

		req.Header = req.Header.WithKVAppend(key, value...)
		return func() {}
	}
}

func OAppendQuery(key string, val string) Option {
	return OAppendQueryH(Q{key: val})
}

func OAppendQueryH(queries IQuery) Option {
	return func(req *RequestBuilder) func() {
		if req.Query == nil {
			req.Query = make(Qs)
		}

		req.Query.Append(queries)
		return func() {}
	}
}
