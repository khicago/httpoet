package httpoet

import (
	"encoding/json"
	"github.com/khicago/irr"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var requestClient = &http.Client{
	Transport: &http.Transport{},
	Timeout:   time.Second * 30,
}

func BuildNRun(req *RequestBuilder, options ...Option) ([]byte, error) {
	if len(options) > 0 {
		defer options[0](req)()
		return BuildNRun(req, options[1:]...)
	}
	req.Build()
	if req.Build().Error != nil {
		return nil, irr.Track(req.Error, "http client do request failed")
	}
	return req.Do()
}

func BuildNRunNParse(req *RequestBuilder, result interface{}, options ...Option) error {
	body, err := BuildNRun(req, options...)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, result); err != nil {
		return irr.Track(err, "unmarshal to result error, type= %v, err= %v", reflect.TypeOf(result), err)
	}
	return nil
}

func (hp *Poet) CreateAbsoluteUrl(url string) string {
	lenH, lenU := len(hp.host), len(url)
	if lenH <= 0 || strings.Contains(url, hp.host) {
		return url
	}
	if lenU <= 0 {
		return hp.host
	}
	if url[0] == '/' && hp.host[lenH-1] == '/' {
		return hp.host + url[1:]
	}
	if url[0] != '/' && hp.host[lenH-1] != '/' {
		return hp.host + "/" + url
	}
	return hp.host + url
}

func (hp *Poet) SpawnReq(url string) *RequestBuilder {
	absoluteUrl := hp.CreateAbsoluteUrl(url)
	return NewReq().XUrl(absoluteUrl).XHeader(hp.baseH)
}

func (hp *Poet) Send(method string, url string, data D, result interface{}, options ...Option) error {
	rb := hp.SpawnReq(url).XMethod(method).XData(data)
	return BuildNRunNParse(rb, result, options...)
}

func (hp *Poet) Get(url string, result interface{}, options ...Option) error {
	rb := hp.SpawnReq(url).XMethod(http.MethodGet)
	return BuildNRunNParse(rb, result, options...)
}

func (hp *Poet) Post(url string, data D, result interface{}, options ...Option) error {
	rb := hp.SpawnReq(url).XMethod(http.MethodPost).XData(data)
	return BuildNRunNParse(rb, result, options...)
}

func (hp *Poet) Put(url string, data D, result interface{}, options ...Option) error {
	rb := hp.SpawnReq(url).XMethod(http.MethodPut).XData(data)
	return BuildNRunNParse(rb, result, options...)
}

func (hp *Poet) Patch(url string, data D, result interface{}, options ...Option) error {
	rb := hp.SpawnReq(url).XMethod(http.MethodPatch).XData(data)
	return BuildNRunNParse(rb, result, options...)
}

func (hp *Poet) Delete(url string, result interface{}, options ...Option) error {
	rb := hp.SpawnReq(url).XMethod(http.MethodDelete)
	return BuildNRunNParse(rb, result, options...)
}
