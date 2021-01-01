package httpoet

import "net/url"

type Q url.Values

func (q Q) Build() string {
	return url.Values(q).Encode()
}

func (q Q) Append(query Q) Q {
	for key, value := range query {
		q[key] = append(q[key], value...)
	}
	return q
}

func (q Q) Add(key string, value string) Q {
	url.Values(q).Add(key, value)
	return q
}

func (q Q) WriteTo(u *url.URL) error {
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
