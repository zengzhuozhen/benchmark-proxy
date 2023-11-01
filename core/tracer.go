package core

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"time"
)

type HttpTracer struct {
	DNSStart             time.Time
	DNSDone              time.Time
	ConnectStart         time.Time
	ConnectDone          time.Time
	GotConn              time.Time
	GotFirstResponseByte time.Time
	TLSHandShakeStart    time.Time
	TLSHandSHakeDone     time.Time
}

type HttpTracerResult struct {
	IsSuccess       bool
	ResponseMessage string
	RequestDataLen  int64
	ResponseDataLen int64
	Duration        DurationInfo
}

type DurationInfo struct {
	Total            int64
	DNSLookup        int64
	TCPConnection    int64
	TLSHandshake     int64
	ServerProcessing int64
	ContentTransfer  int64
}

func NewDurationInfo(t *HttpTracer) DurationInfo {
	now := time.Now()
	dnsLookup := t.DNSDone.Sub(t.DNSStart).Milliseconds()
	tcpConnection := t.ConnectDone.Sub(t.ConnectStart).Milliseconds()
	tlsHandshake := t.TLSHandSHakeDone.Sub(t.TLSHandShakeStart).Milliseconds()
	ServerProcessing := t.GotFirstResponseByte.Sub(t.GotConn).Milliseconds()

	start := t.GotConn
	if !t.ConnectStart.IsZero() {
		start = t.ConnectStart
	}
	if !t.DNSStart.IsZero() {
		start = t.DNSStart
	}

	return DurationInfo{
		Total:            now.Sub(start).Milliseconds(),
		DNSLookup:        dnsLookup,
		TCPConnection:    tcpConnection,
		TLSHandshake:     tlsHandshake,
		ServerProcessing: ServerProcessing,
		ContentTransfer:  now.Sub(t.GotFirstResponseByte).Milliseconds(),
	}
}

func (t *HttpTracer) Trace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			t.DNSStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			t.DNSDone = time.Now()
		},
		ConnectStart: func(network, addr string) {
			t.ConnectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				fmt.Println("http trace error.ConnectDone:", err.Error())
			}
			t.ConnectDone = time.Now()
		},
		GotConn: func(info httptrace.GotConnInfo) {
			t.GotConn = time.Now()
		},
		GotFirstResponseByte: func() {
			t.GotFirstResponseByte = time.Now()
		},
		TLSHandshakeStart: func() {
			t.TLSHandShakeStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			if err != nil {
				fmt.Println("http trace error.TLSHandshakeDone:", err.Error())
			}
			t.TLSHandSHakeDone = time.Now()
		},
	}
}

func (t *HttpTracer) Result(req *http.Request, resp *http.Response, checker *ResponseChecker) HttpTracerResult {
	result := new(HttpTracerResult)
	result.RequestDataLen = req.ContentLength
	result.ResponseDataLen = resp.ContentLength
	msg, _ := ioutil.ReadAll(resp.Body)
	result.ResponseMessage = string(msg)
	result.IsSuccess = checker.Check(resp.StatusCode, result.ResponseMessage)
	result.Duration = NewDurationInfo(t)
	return *result
}
