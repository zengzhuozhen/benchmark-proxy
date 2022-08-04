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
	port    int
	rootCA  *x509.Certificate
	rootKey *rsa.PrivateKey
}

func NewBenchProxyService(port int, rootCA *x509.Certificate, rootKey *rsa.PrivateKey) *BenchmarkProxyService {
	return &BenchmarkProxyService{port: port, rootCA: rootCA, rootKey: rootKey}
}

func (s *BenchmarkProxyService) Serve() {
	_ = http.ListenAndServe(fmt.Sprintf(":%d", s.port), s)
}

func (s *BenchmarkProxyService) ServeHTTP(originRespWriter http.ResponseWriter, originReq *http.Request) {
	s.WrapInTls(originReq, originRespWriter, func(req *http.Request, respWriter http.ResponseWriter) {
		executor := NewExecutor(req)
		executor.Run()
		respWriter.WriteHeader(http.StatusOK)
		respWriter.Write(executor.Result().Print())
	})
}

func (s *BenchmarkProxyService) WrapInTls(originReq *http.Request, originRespWriter http.ResponseWriter, fn func(originReq *http.Request, originRespWriter http.ResponseWriter)) {
	if originReq.Method != http.MethodConnect {
		fn(originReq, originRespWriter)
		return
	}
	hijacker, ok := originRespWriter.(http.Hijacker)
	if !ok {
		http.Error(originRespWriter, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(originRespWriter, fmt.Sprintf("Hijacking failed: %s", err), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()
	clientConn.Write([]byte("HTTP/1.1 200 Connection established\n\n"))

	tlsConfig, err := GenerateTlsConfig(originReq.URL.Host, s.rootCA, s.rootKey)
	if err != nil {
		http.Error(originRespWriter, fmt.Sprintf("HTTPS生成证书失败: %s", err), http.StatusServiceUnavailable)
		return
	}
	tlsClientConn := tls.Server(clientConn, tlsConfig)
	defer tlsClientConn.Close()
	if err := tlsClientConn.Handshake(); err != nil {
		http.Error(originRespWriter, fmt.Sprintf("HTTPS解密, 握手失败: %s", err), http.StatusServiceUnavailable)
		return
	}

	buf := bufio.NewReader(tlsClientConn)
	tlsReq, err := http.ReadRequest(buf)
	if err != nil {
		if err != io.EOF {
			http.Error(originRespWriter, fmt.Sprintf("HTTPS解密, 读取客户端请求失败:%s", err), http.StatusServiceUnavailable)
		}
		return
	}
	tlsReq.RemoteAddr = originReq.RemoteAddr
	tlsReq.URL.Scheme = "https"
	tlsReq.URL.Host = tlsReq.Host

	if !strings.Contains(tlsReq.URL.Host, ":") {
		tlsReq.URL.Host += ":443"
	}

	targetConn, err := net.DialTimeout("tcp", tlsReq.URL.Host, 5*time.Second)
	if err != nil {
		http.Error(originRespWriter, fmt.Sprintf("隧道转发连接目标服务器失败, :%s", err), http.StatusServiceUnavailable)
		return
	}
	defer targetConn.Close()
	targetConn.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", originReq.Host)))

	fn(tlsReq, &BenchmarkRespWriter{tlsClientConn})
}
