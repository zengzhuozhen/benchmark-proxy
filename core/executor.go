package core

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zengzhuozhen/benchmark-proxy/protocol"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"time"
)

const (
	DebugRequestMessageFormat  = "Request Message : %+v \n"
	DebugRequestErrorFormat    = "Request Error : %s \n"
	DebugResponseMessageFormat = "Response Message : %s \n"
)

type BenchmarkExecutorType int

type BenchmarkProxyHeader struct {
	ExecTimes       int
	ExecDuration    time.Duration
	ExecConcurrency int
	ResponseChecker *ResponseChecker
}

type BenchmarkReqConfig struct {
	proxyHeaders *BenchmarkProxyHeader
	body         []byte
	originReq    *http.Request
}

type BenchmarkExecutor interface {
	Run() error // using goroutine pool
	ClearHopHeaders(header http.Header) http.Header
	Result() *Statistic
	MetaData() *BenchmarkProxyHeader
}

func NewExecutor(req *http.Request) BenchmarkExecutor {
	header := make(http.Header)
	protocol.CopyHeader(header, req.Header)
	body, _ := ioutil.ReadAll(req.Body)
	for _, i := range protocol.BenchmarkProxyHeaders {
		req.Header.Del(i)
	}
	proxyHeader := NewProxyHeader(header)
	if proxyHeader.ExecTimes == 0 && proxyHeader.ExecDuration == time.Duration(0) {
		panic(fmt.Sprintf("unKnow Benchmark Proxy type, use one of these (%s and %s) header ",
			protocol.BenchmarkProxyExecTimes, protocol.BenchmarkProxyExecDuration))
	}
	if proxyHeader.ExecTimes > 0 {
		return &BenchmarkExecTimes{
			ctx: context.Background(),
			BenchmarkReqConfig: BenchmarkReqConfig{
				proxyHeaders: proxyHeader,
				body:         body,
				originReq:    req,
			},
			Executor: Executor{new(Statistic), make(chan HttpTracerResult)},
		}
	}
	return &BenchmarkExecDuration{
		ctx: context.Background(),
		BenchmarkReqConfig: BenchmarkReqConfig{
			proxyHeaders: proxyHeader,
			body:         body,
			originReq:    req,
		},
		Executor: Executor{new(Statistic), make(chan HttpTracerResult)},
	}

}

func NewProxyHeader(header http.Header) *BenchmarkProxyHeader {
	var execTime, execDuration, execConcurrency, checkResultStatus int
	execTimesStr := protocol.GetProxyHeaderParam(header, protocol.BenchmarkProxyExecTimes)
	if execTimesStr != "" {
		execTime, _ = strconv.Atoi(execTimesStr)
	}

	execDurationStr := protocol.GetProxyHeaderParam(header, protocol.BenchmarkProxyExecDuration)
	if execDurationStr != "" {
		execDuration, _ = strconv.Atoi(execDurationStr)
	}

	execConcurrencyStr := protocol.GetProxyHeaderParam(header, protocol.BenchmarkProxyConcurrency)
	if execConcurrencyStr != "" {
		execConcurrency, _ = strconv.Atoi(execConcurrencyStr)
	}

	var responseCheckerOptions []ResponseCheckOption

	checkResultStatusStr := protocol.GetProxyHeaderParam(header, protocol.BenchmarkProxyCheckResultStatus)
	if checkResultStatusStr != "" {
		checkResultStatus, _ = strconv.Atoi(checkResultStatusStr)
		responseCheckerOptions = append(responseCheckerOptions, ResponseCheckerStatusRule(checkResultStatus))
	}

	checkResultBody := protocol.GetProxyHeaderParam(header, protocol.BenchmarkProxyCheckResultBody)
	if checkResultBody != "" {
		responseCheckerOptions = append(responseCheckerOptions, ResponseCheckerBodyRule(checkResultBody))
	}

	responseChecker := NewResponseChecker(responseCheckerOptions...)

	return &BenchmarkProxyHeader{
		ExecTimes:       execTime,
		ExecDuration:    time.Duration(execDuration * int(time.Second)),
		ExecConcurrency: execConcurrency,
		ResponseChecker: responseChecker,
	}
}

type BenchmarkExecTimes struct {
	ctx context.Context
	BenchmarkReqConfig
	Executor
}

func (exec *BenchmarkExecTimes) MetaData() *BenchmarkProxyHeader {
	return exec.BenchmarkReqConfig.proxyHeaders
}

type BenchmarkExecDuration struct {
	ctx context.Context
	BenchmarkReqConfig
	Executor
}

func (exec *BenchmarkExecDuration) MetaData() *BenchmarkProxyHeader {
	return exec.BenchmarkReqConfig.proxyHeaders
}

func (config *BenchmarkReqConfig) ClearHopHeaders(originHeader http.Header) http.Header {
	for _, i := range protocol.HopHeaders {
		originHeader.Del(i)
	}
	return originHeader
}

func (exec *BenchmarkExecTimes) Run() error {
	var err error
	urlParser, bodyParser := NewTagCompoundParser(), NewTagCompoundParser()
	pool.setCap(exec.proxyHeaders.ExecConcurrency)
	go exec.statistic.Aggregate(exec.resultChan)
	for i := 0; i < exec.proxyHeaders.ExecTimes; i++ {
		task := createTask(func() {
			newReq := exec.originReq.Clone(exec.originReq.Context())
			newReq.Body = io.NopCloser(bytes.NewReader(exec.body))
			newReq.Header = exec.ClearHopHeaders(newReq.Header)
			exec.ReplaceCustomizeTag(urlParser, bodyParser, newReq)
			result, reqErr := exec.RunOnce(newReq, exec.proxyHeaders.ResponseChecker)
			exec.resultChan <- result
			if reqErr != nil {
				err = reqErr
				log.Debugf(DebugRequestErrorFormat, err.Error())
			}
		})
		pool.addTask(task)
	}
	pool.wait()
	close(exec.resultChan)
	return err
}

func (exec *BenchmarkExecDuration) Run() error {
	var err error
	pool.setCap(exec.proxyHeaders.ExecConcurrency)
	ctx := exec.originReq.Context()
	childCtx, cancelFunc := context.WithCancel(ctx)
	go time.AfterFunc(exec.proxyHeaders.ExecDuration, cancelFunc)
	go exec.statistic.Aggregate(exec.resultChan)
	urlParser, bodyParser := NewTagCompoundParser(), NewTagCompoundParser()
	for {
		select {
		case <-childCtx.Done():
			goto OUT
		default:
			task := createTask(func() {
				newReq := exec.originReq.Clone(exec.originReq.Context())
				newReq.Body = io.NopCloser(bytes.NewReader(exec.body))
				newReq.Header = exec.ClearHopHeaders(newReq.Header)
				exec.ReplaceCustomizeTag(urlParser, bodyParser, newReq)
				result, reqErr := exec.RunOnce(newReq, exec.proxyHeaders.ResponseChecker)
				exec.resultChan <- result
				if reqErr != nil {
					err = reqErr
					log.Debugf(DebugRequestErrorFormat, reqErr.Error())
				}
			})
			pool.addTask(task)
		}
	}
OUT:
	pool.wait()
	close(exec.resultChan)
	return err
}

type Executor struct {
	statistic  *Statistic
	resultChan chan HttpTracerResult
}

func (exec *Executor) ReplaceCustomizeTag(urlParser, bodyParser *TagCompoundParser, req *http.Request) {
	// url
	queryPairs := req.URL.Query()
	for k, v := range queryPairs {
		// 1. get need replace tag
		var replace []string
		for _, i := range v {
			replace = append(replace, urlParser.ParseCustomizeTag(i)) // 2.replace every tag
		}
		queryPairs.Del(k)
		for _, i := range replace { // 3.reset queryParis
			queryPairs.Add(k, i)
		}
	}
	req.URL.RawQuery = queryPairs.Encode()
	// body
	bodyContent, _ := ioutil.ReadAll(req.Body)
	defer func() {
		req.Body = io.NopCloser(bytes.NewReader(bodyContent))
	}()
	parseContent := bodyParser.ParseCustomizeTag(string(bodyContent))
	bodyContent = []byte(parseContent)
}

func (exec *Executor) RunOnce(req *http.Request, checker *ResponseChecker) (HttpTracerResult, error) {
	tracer := &HttpTracer{}
	execReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	execReq.Header = req.Header
	if err != nil {
		return HttpTracerResult{}, fmt.Errorf("构造请求失败,错误原因:%s", err.Error())
	}
	log.Debugf(DebugRequestMessageFormat, execReq)
	execReq = execReq.WithContext(httptrace.WithClientTrace(execReq.Context(), tracer.Trace()))
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	resp, err := client.Do(execReq)
	if err != nil {
		return HttpTracerResult{}, fmt.Errorf("执行请求失败,错误原因:%s", err.Error())
	}
	defer resp.Body.Close()
	result := tracer.Result(execReq, resp, checker)
	log.Debugf(DebugResponseMessageFormat, result.ResponseMessage)
	return result, nil
}

func (exec *Executor) Result() *Statistic {
	return exec.statistic
}
