package protocol

import (
	"net/http"
)

// HopHeaders Hop-by-hop proxyHeaders. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop proxyHeaders are required to appear in the
// Connection header field. These are the proxyHeaders defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
var HopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

var BenchmarkProxyHeaders = []string{
	BenchmarkProxyExecTimes,
	BenchmarkProxyExecDuration,
	BenchmarkProxyConcurrency,
	BenchmarkProxyCheckResultStatus,
	BenchmarkProxyCheckResultBody,
	BenchmarkProxyAuth,
}

const (
	BenchmarkProxyExecTimes         = "Benchmark-Proxy-Times"
	BenchmarkProxyExecDuration      = "Benchmark-Proxy-Duration"
	BenchmarkProxyConcurrency       = "Benchmark-Proxy-Concurrency"
	BenchmarkProxyCheckResultStatus = "Benchmark-Proxy-Check-Result-Status"
	BenchmarkProxyCheckResultBody   = "Benchmark-Proxy-Check-Result-Body"
	BenchmarkProxyAuth              = "Benchmark-Proxy-Auth"
)

func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func GetProxyHeaderParam(header http.Header, key string) string {
	if val, ok := header[key]; ok {
		return val[0]
	}
	return ""
}
