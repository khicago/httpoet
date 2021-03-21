package httpoet

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/khicago/irr"
)

type RequestBuilder struct {
	Method string
	Url    string
	Header IHeader
	Query  IQuery
	Data   interface{}

	Context context.Context

	Cookies []*http.Cookie

	req   *http.Request
	Error error
}

func NewReq() *RequestBuilder {
	return &RequestBuilder{}
}

func (rb *RequestBuilder) ResetErrorState() *RequestBuilder {
	rb.Error = nil
	return rb
}

func (rb *RequestBuilder) XContext(ctx context.Context) *RequestBuilder {
	if rb.Error != nil {
		return rb
	}
	rb.Context = ctx
	return rb
}

func (rb *RequestBuilder) XMethod(m string) *RequestBuilder {
	if rb.Error != nil {
		return rb
	}
	rb.Method = m
	return rb
}

func (rb *RequestBuilder) XUrl(u string) *RequestBuilder {
	if rb.Error != nil {
		return rb
	}
	rb.Url = u
	return rb
}

func (rb *RequestBuilder) XHeader(h IHeader) *RequestBuilder {
	if rb.Error != nil {
		return rb
	}
	rb.Header = h
	return rb
}

func (rb *RequestBuilder) XQuery(q IQuery) *RequestBuilder {
	if rb.Error != nil {
		return rb
	}
	rb.Query = q
	return rb
}

func (rb *RequestBuilder) XData(d D) *RequestBuilder {
	return rb.XDataCustom(d)
}

func (rb *RequestBuilder) XDataCustom(d interface{}) *RequestBuilder {
	if rb.Error != nil {
		return rb
	}
	rb.Data = d
	return rb
}

func (rb *RequestBuilder) Build() (rbRet *RequestBuilder) {
	if rb.Error != nil {
		return rb
	}
	rbRet = rb

	u, err := url.Parse(rb.Url)
	if err != nil {
		rb.Error = irr.Track(err, "build url= %s failed", rb.Url)
		return
	}

	if rb.Query != nil {
		if err = rb.Query.WriteTo(u); err != nil {
			rb.Error = irr.Track(err, "build query= %s failed", rb.Url)
			return
		}
	}

	byteArr, err := rb.buildData()
	if err != nil {
		rb.Error = irr.Track(err, "build data= %v failed", rb.Data)
		return
	}

	req, err := http.NewRequest(rb.Method, u.String(), bytes.NewReader(byteArr))
	if err != nil {
		rb.Error = irr.Track(err, "create http request failed")
		return
	}

	if rb.Header != nil {
		rb.Header.foreach(func(k string, vs ...string) {
			req.Header[k] = append(req.Header[k], vs...)
		})
	}

	if rb.Cookies != nil {
		for _, c := range rb.Cookies {
			req.AddCookie(c)
		}

		fmt.Println("build cookie", rb.Cookies, req.Header)
	}

	rb.req = req
	return
}

func (rb *RequestBuilder) Do() ([]byte, error) {
	if rb.Error != nil {
		return nil, irr.Track(rb.Error, "cannot do incorrect request")
	}
	resp, err := RequestClient.Do(rb.req)
	if err != nil {
		return nil, irr.Track(err, "http client do request failed")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, irr.Track(err, "read buffer error")
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, irr.Track(err, "status error, status= %d, body= %s", resp.StatusCode, body)
	}
	return body, nil
}

func (rb *RequestBuilder) Request() *http.Request {
	return rb.req
}

func (rb *RequestBuilder) buildData() ([]byte, error) {
	switch t := rb.Data.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
		return json.Marshal(rb.Data)
	}
}
