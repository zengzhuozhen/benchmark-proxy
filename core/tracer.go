package core

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"
)

type HttpTracer struct {
	GetConnTime              int64
	GotConnTime              int64
	ConnectStartTime         int64
	ConnectDoneTime          int64
	TLSHandshakeStartTime    int64
	TLSHandshakeDoneTime     int64
	WroteRequestTime         int64
	GotFirstResponseByteTime int64
}

type HttpTracerResult struct {
	IsSuccess       bool
	ResponseMessage string
	RequestDataLen  int64
	ResponseDataLen int64
	Duration        time.Duration
}

func (t *HttpTracer) Trace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			t.GetConnTime = time.Now().UnixNano()
		},
		ConnectStart: func(network, addr string) {
			atomic.CompareAndSwapInt64(&t.ConnectStartTime, 0, time.Now().UnixNano())
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				atomic.CompareAndSwapInt64(&t.ConnectDoneTime, 0, time.Now().UnixNano())
			}
		},
		TLSHandshakeStart: func() {
			atomic.CompareAndSwapInt64(&t.TLSHandshakeStartTime, 0, time.Now().UnixNano())
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			if err != nil {
				atomic.CompareAndSwapInt64(&t.TLSHandshakeDoneTime, 0, time.Now().UnixNano())
			}
		},
		GotConn: func(info httptrace.GotConnInfo) {
			atomic.CompareAndSwapInt64(&t.GotConnTime, 0, time.Now().UnixNano())
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			if info.Err != nil {
				atomic.StoreInt64(&t.WroteRequestTime, time.Now().UnixNano())
			}
		},
		GotFirstResponseByte: func() {
			atomic.CompareAndSwapInt64(&t.GotFirstResponseByteTime, 0, time.Now().UnixNano())
		},
	}
}

func (t *HttpTracer) Result(req *http.Request, resp *http.Response) HttpTracerResult {
	result := new(HttpTracerResult)
	result.IsSuccess = resp.StatusCode == http.StatusOK
	result.RequestDataLen = req.ContentLength
	result.ResponseDataLen = resp.ContentLength
	_, _ = resp.Body.Read([]byte(result.ResponseMessage))
	result.Duration = time.Duration(t.GotFirstResponseByteTime - t.GetConnTime)
	return *result
}
