package httpoet

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type (
	queryIdentity struct{}

	IQuery interface {
		foreach(fn func(key string, value ...string)) *queryIdentity

		Build() string
		Append(query IQuery) IQuery
		WriteTo(u *url.URL) error
		WriteToPth(u string, args ...interface{}) (string, error)
	}

	Q  map[string]string
	Qs map[string][]string
)

var _, _ IQuery = Q{}, Qs{}

///////////////// Region q

func (q Q) foreach(fn func(key string, value ...string)) *queryIdentity {
	for k, v := range q {
		fn(k, v)
	}
	return &queryIdentity{}
}

func (q Q) Build() string {
	if q == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(q[k]))
	}
	return buf.String()
}

func (q Q) Set(key string, value string) IQuery {
	q[key] = value
	return q
}

func (q Q) Append(query IQuery) IQuery {
	newQ := make(Qs)
	for k, v := range q {
		newQ[k] = []string{v}
	}
	return newQ.Append(query)
}

func (q Q) WriteTo(u *url.URL) error {
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return err
	}
	for key, value := range q {
		query[key] = append(query[key], value)
	}
	u.RawQuery = query.Encode()
	return nil
}

func (q Q) WriteToPth(u string, args ...interface{}) (string, error) {
	if len(args) > 0 {
		u = fmt.Sprintf(u, args...)
	}
	tu, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	err = q.WriteTo(tu)
	if err != nil {
		return "", err
	}
	return tu.String(), nil
}

///////////////// Region qs

func (q Qs) foreach(fn func(key string, value ...string)) *queryIdentity {
	for k, vs := range q {
		fn(k, vs...)
	}
	return &queryIdentity{}
}

func (q Qs) Build() string {
	return url.Values(q).Encode()
}

func (q Qs) Append(query IQuery) IQuery {
	query.foreach(func(key string, vs ...string) {
		q[key] = append(q[key], vs...)
	})
	return q
}

func (q Qs) Add(key string, value string) Qs {
	url.Values(q).Add(key, value)
	return q
}

func (q Qs) WriteTo(u *url.URL) error {
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return err
	}
	for key, value := range q {
		query[key] = append(query[key], value...)
	}
	u.RawQuery = query.Encode()
	return nil
}

func (q Qs) WriteToPth(u string, args ...interface{}) (string, error) {
	if len(args) > 0 {
		u = fmt.Sprintf(u, args...)
	}
	tu, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	err = q.WriteTo(tu)
	if err != nil {
		return "", err
	}
	return tu.String(), nil
}
