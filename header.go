package httpoet

import (
	"net/http"
)

type H http.Header

func (h H) Append(headers H) H {
	for key, value := range headers {
		h.Add(key, value...)
	}
	return h
}

func (h H) Add(key string, value ...string) H {
	h[key] = append(h[key], value...)
	return h
}

func (h H) Set(key string, value ...string) H {
	h[key] = value
	return h
}
