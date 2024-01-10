package core

import (
	log "github.com/sirupsen/logrus"
	"github.com/zengzhuozhen/benchmark-proxy/protocol"
	"github.com/zengzhuozhen/benchmark-proxy/resources"
	"net/http"
)

type Delegate interface {
	Auth(req *http.Request) bool
	Handle(req *http.Request, resp http.ResponseWriter)
	Response(resp http.ResponseWriter, code int, msg []byte)
}

var _ Delegate = &DefaultDelegate{}

type DefaultDelegate struct {
	Delegate
}

func (d *DefaultDelegate) Auth(req *http.Request) bool {
	key := req.Header.Get(protocol.BenchmarkProxyAuth)
	if resources.AuthKey != "" && resources.AuthKey != key {
		return false
	}
	return true
}

func (d *DefaultDelegate) Handle(req *http.Request, resp http.ResponseWriter) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	if d.Auth(req) == false {
		d.Response(resp, http.StatusUnauthorized, []byte("请求未授权"))
		return
	}
	executor := NewExecutor(req)
	if err := executor.Run(); err != nil {
		d.Response(resp, http.StatusBadRequest, []byte(err.Error()))
		return
	}
	d.Response(resp, http.StatusOK, executor.Result().Print(executor.MetaData()))
}

func (d *DefaultDelegate) Response(resp http.ResponseWriter, code int, msg []byte) {
	resp.WriteHeader(code)
	resp.Write(msg)
}
