package httpoet

import (
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/khicago/irr"
)

var RequestClient = &http.Client{
	Transport: &http.Transport{},
	Timeout:   time.Second * 30,
}

type (
	result struct {
		body []byte
		err  error
	}

	IResult interface {
		Body() ([]byte, error)
		ParseJson(resultObject interface{}) error
		ParseCustom(resultObject interface{}, method func(body []byte, target interface{}) error) error
	}
)

func (rr *result) Body() ([]byte, error) {
	return rr.body, rr.err
}

func (rr *result) ParseJson(resultObject interface{}) error {
	if rr.err != nil {
		return irr.TrackSkip(1, rr.err, "parse stopped, error already exist")
	}
	if rr.err = json.Unmarshal(rr.body, resultObject); rr.err != nil {
		return irr.TrackSkip(1, rr.err, "unmarshal to result error, type= %v", reflect.TypeOf(resultObject))
	}
	return rr.err
}

func (rr *result) ParseCustom(resultObject interface{}, method func(body []byte, target interface{}) error) error {
	if rr.err != nil {
		return irr.TrackSkip(1, rr.err, "parse stopped, error already exist")
	}
	if method == nil {
		return irr.TraceSkip(1, "custom parser cannot be empty")
	}
	if rr.err = method(rr.body, resultObject); rr.err != nil {
		return irr.TrackSkip(1, rr.err, "unmarshal to result error, type= %v", reflect.TypeOf(resultObject))
	}
	return rr.err
}

func BuildNRun(req *RequestBuilder, options ...Option) IResult {
	if len(options) > 0 {
		defer options[0](req)()
		return BuildNRun(req, options[1:]...)
	}

	if req.Build().Error != nil {
		return &result{err: irr.Track(req.Error, "http client do request failed")}
	}

	// fmt.Printf("req %+v", req)
	body, err := req.Do()
	return &result{body, err}
}

func (hp *Poet) GetHost() string {
	n := len(hp.hosts)
	switch n {
	case 0:
		return "localhost"
	case 1:
		return hp.hosts[0]
	default:
		ind := rand.Intn(len(hp.hosts))
		return hp.hosts[ind]
	}
}

func (hp *Poet) CreateAbsoluteUrl(url string) string {
	host := hp.GetHost()

	lenH, lenU := len(host), len(url)
	if lenH <= 0 || strings.Contains(url, host) {
		return url
	}
	if lenU <= 0 {
		return host
	}
	if url[0] == '/' && host[lenH-1] == '/' {
		return host + url[1:]
	}
	if url[0] != '/' && host[lenH-1] != '/' {
		return host + "/" + url
	}
	return host + url
}

func (hp *Poet) SpawnReq(url string) *RequestBuilder {
	absoluteUrl := hp.CreateAbsoluteUrl(url)
	return NewReq().XUrl(absoluteUrl).XHeader(hp.baseH)
}

func (hp *Poet) Send(method string, url string, data interface{}, options ...Option) IResult {
	rb := hp.SpawnReq(url).XMethod(method).XDataCustom(data)
	return BuildNRun(rb, options...)
}

/* REST-ful apis */

func (hp *Poet) Post(url string, data interface{}, options ...Option) IResult {
	rb := hp.SpawnReq(url).XMethod(http.MethodPost).XDataCustom(data)
	return BuildNRun(rb, options...)
}

func (hp *Poet) Put(url string, data interface{}, options ...Option) IResult {
	rb := hp.SpawnReq(url).XMethod(http.MethodPut).XDataCustom(data)
	return BuildNRun(rb, options...)
}

func (hp *Poet) Patch(url string, data interface{}, options ...Option) IResult {
	rb := hp.SpawnReq(url).XMethod(http.MethodPatch).XDataCustom(data)
	return BuildNRun(rb, options...)
}

func (hp *Poet) Get(url string, options ...Option) IResult {
	rb := hp.SpawnReq(url).XMethod(http.MethodGet)
	return BuildNRun(rb, options...)
}

func (hp *Poet) Delete(url string, options ...Option) IResult {
	rb := hp.SpawnReq(url).XMethod(http.MethodDelete)
	return BuildNRun(rb, options...)
}
