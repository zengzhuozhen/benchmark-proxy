package core

import (
	"bufio"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type BenchmarkProxyService struct {
	delegate Delegate
	port     int
	rootCA   *x509.Certificate
	rootKey  *rsa.PrivateKey
}

func NewBenchProxyService(port int, rootCA *x509.Certificate, rootKey *rsa.PrivateKey) *BenchmarkProxyService {
	return &BenchmarkProxyService{port: port, rootCA: rootCA, rootKey: rootKey, delegate: &DefaultDelegate{}}
}

func (s *BenchmarkProxyService) Serve() {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s); err != nil {
		panic(err)
	}
}

func (s *BenchmarkProxyService) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodConnect:
		s.tunnelProxy(req, resp, s.delegate.Handle)
	default:
		s.httpProxy(req, resp, s.delegate.Handle)
	}
}

func (s *BenchmarkProxyService) httpProxy(req *http.Request, resp http.ResponseWriter, fn func(req *http.Request, resp http.ResponseWriter)) {
	req.URL.Scheme = "http"
	req.URL.Host = req.Host
	fn(req, resp)
}

func (s *BenchmarkProxyService) tunnelProxy(req *http.Request, resp http.ResponseWriter, fn func(req *http.Request, resp http.ResponseWriter)) {
	if s.rootKey == nil || s.rootCA == nil {
		http.Error(resp, "未加载证书，不支持https协议", http.StatusBadRequest)
		return
	}
	hijacker, ok := resp.(http.Hijacker)
	if !ok {
		http.Error(resp, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(resp, fmt.Sprintf("Hijacking failed: %s", err), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()
	clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

	tlsConfig, err := GenerateTlsConfig(req.URL.Host, s.rootCA, s.rootKey)
	if err != nil {
		http.Error(resp, fmt.Sprintf("HTTPS生成证书失败: %s", err), http.StatusServiceUnavailable)
		return
	}
	tlsClientConn := tls.Server(clientConn, tlsConfig)
	defer tlsClientConn.Close()
	if err := tlsClientConn.Handshake(); err != nil {
		http.Error(resp, fmt.Sprintf("HTTPS解密, 握手失败: %s", err), http.StatusServiceUnavailable)
		return
	}

	buf := bufio.NewReader(tlsClientConn)
	tlsReq, err := http.ReadRequest(buf)
	if err != nil {
		if err != io.EOF {
			http.Error(resp, fmt.Sprintf("HTTPS解密, 读取客户端请求失败:%s", err), http.StatusServiceUnavailable)
		}
		return
	}
	tlsReq.RemoteAddr = req.RemoteAddr
	tlsReq.URL.Scheme = "https"
	tlsReq.URL.Host = tlsReq.Host

	if !strings.Contains(tlsReq.URL.Host, ":") {
		tlsReq.URL.Host += ":443"
	}

	targetConn, err := net.DialTimeout("tcp", tlsReq.URL.Host, 5*time.Second)
	if err != nil {
		http.Error(resp, fmt.Sprintf("隧道转发连接目标服务器失败, :%s", err), http.StatusServiceUnavailable)
		return
	}
	defer targetConn.Close()
	targetConn.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", req.Host)))

	fn(tlsReq, &BenchmarkRespWriter{tlsClientConn})
}
