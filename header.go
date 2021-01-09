package httpoet

import (
	"net/http"
)

type H http.Header

func (h H) WithH(headers H) H {
	newH := make(H)
	for k, v := range h {
		newH[k] = v
	}
	for key, value := range headers {
		newH[key] = append(newH[key], value...)
	}
	return h
}

func (h H) WithKV(key string, value ...string) H {
	newH := make(H)
	for k, v := range h {
		newH[k] = v
	}
	newH[key] = value
	return newH
}
