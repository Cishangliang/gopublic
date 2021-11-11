package requester

import (
	"net/http"
	"time"

	"github.com/imroc/req"
)

type requester struct {
	r *req.Req
}

var instance *requester

func init() {
	trans, _ := req.Client().Transport.(*http.Transport)
	trans.MaxIdleConns = 50
	trans.TLSHandshakeTimeout = 5 * time.Second
	trans.DisableKeepAlives = true
	instance = &requester{
		r: req.New(),
	}
}

func Instance() *requester {
	return instance
}

func (r *requester) Get(url string) (*req.Resp, error) {
	return r.r.Get(url)
}

func (r *requester) GetWithParam(url string, param req.Param) (*req.Resp, error) {
	return r.r.Get(url, param)
}

func (r *requester) Post(url string, body interface{}) (*req.Resp, error) {
	return r.r.Post(url, req.BodyJSON(body))
}

func (r *requester) GetWithQuery(url string, param map[string]string) (*req.Resp, error) {
	var p = req.Param{}
	for k, v := range param {
		p[k] = v
	}
	return r.r.Get(url, p)
}

func (r *requester) Put(url string, body interface{}) (*req.Resp, error) {
	return r.r.Put(url, body)
}

func (r *requester) Delete(url string, body interface{}) (*req.Resp, error) {
	return r.r.Delete(url, body)
}
