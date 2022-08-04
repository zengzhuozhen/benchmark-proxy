package core

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type BenchmarkRespWriter struct {
	*tls.Conn
}

func (t *BenchmarkRespWriter) Header() http.Header {
	emptyHeader := make(http.Header)
	return emptyHeader
}

func (t *BenchmarkRespWriter) WriteHeader(statusCode int) {
	t.Write([]byte(fmt.Sprintf("HTTP/1.1 %d \n\n", statusCode)))
}
