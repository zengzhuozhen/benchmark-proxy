package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/zengzhuozhen/benchmark-proxy/protocol"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"sync"
	"time"
)

type BenchmarkExecutorType int

type BenchmarkProxyHeader struct {
	ExecTimes       int
	ExecDuration    time.Duration
	ExecConcurrency int
}

type BenchmarkReqConfig struct {
	proxyHeaders *BenchmarkProxyHeader
	body         []byte
	originReq    *http.Request
}

type BenchmarkExecutor interface {
	Run(isDebug bool)
	ClearHopHeaders(header http.Header) http.Header
	Result() *Statistic
}

func NewExecutor(req *http.Request) BenchmarkExecutor {
	header := make(http.Header)
	protocol.CopyHeader(header, req.Header)
	body, _ := ioutil.ReadAll(req.Body)
	for _, i := range protocol.BenchmarkProxyHeaders {
		req.Header.Del(i)
	}
	proxyHeader := NewProxyHeader(header)
	if protocol.CheckProxyHeader(header, protocol.BenchProxyTypeTime) {
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
	if protocol.CheckProxyHeader(header, protocol.BenchProxyTypeDuration) {
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
	panic(fmt.Sprintf("unKnow Benchmark Proxy type, did you forget set %s Header or with wrong value (%s or %s)?",
		protocol.BenchmarkProxyType, protocol.BenchProxyTypeTime, protocol.BenchProxyTypeDuration))
}

func NewProxyHeader(header http.Header) *BenchmarkProxyHeader {
	var execTime, execDuration, execConcurrency int
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
	return &BenchmarkProxyHeader{
		ExecTimes:       execTime,
		ExecDuration:    time.Duration(execDuration * int(time.Second)),
		ExecConcurrency: execConcurrency,
	}
}

type BenchmarkExecTimes struct {
	ctx context.Context
	BenchmarkReqConfig
	Executor
}

type BenchmarkExecDuration struct {
	ctx context.Context
	BenchmarkReqConfig
	Executor
}

func (config *BenchmarkReqConfig) ClearHopHeaders(originHeader http.Header) http.Header {
	for _, i := range protocol.HopHeaders {
		originHeader.Del(i)
	}
	return originHeader
}

func (exec *BenchmarkExecTimes) Run(isDebug bool) {
	wg := sync.WaitGroup{}
	wg.Add(exec.proxyHeaders.ExecTimes)
	concurrencyBuffer := make(chan struct{}, exec.proxyHeaders.ExecConcurrency)
	go exec.statistic.Aggregate(exec.resultChan)
	urlParser, bodyParser := NewTagCompoundParser(), NewTagCompoundParser()
	for i := 0; i < exec.proxyHeaders.ExecTimes; i++ {
		go func() {
			defer wg.Done()
			concurrencyBuffer <- struct{}{}
			newReq := exec.originReq.Clone(exec.originReq.Context())
			newReq.Body = io.NopCloser(bytes.NewReader(exec.body))
			newReq.Header = exec.ClearHopHeaders(newReq.Header)
			exec.ReplaceCustomizeTag(urlParser, bodyParser, newReq)
			result, err := exec.RunOnce(newReq)
			if isDebug {
				fmt.Printf("「DEBUG」Response Message : %s \n", result.ResponseMessage)
			}
			if err == nil {
				exec.resultChan <- result
			}
			<-concurrencyBuffer
		}()
	}
	wg.Wait()
	close(concurrencyBuffer)
	close(exec.resultChan)
}

func (exec *BenchmarkExecDuration) Run(isDebug bool) {
	ctx := exec.originReq.Context()
	childCtx, cancelFunc := context.WithCancel(ctx)
	concurrencyBuffer := make(chan struct{}, exec.proxyHeaders.ExecConcurrency)
	wg := &sync.WaitGroup{}
	go time.AfterFunc(exec.proxyHeaders.ExecDuration, cancelFunc)
	go exec.statistic.Aggregate(exec.resultChan)
	urlParser, bodyParser := NewTagCompoundParser(), NewTagCompoundParser()
	for {
		select {
		case <-childCtx.Done():
			goto OUT
		default:
			concurrencyBuffer <- struct{}{}
			go func() {
				wg.Add(1)
				defer wg.Done()
				newReq := exec.originReq.Clone(childCtx)
				newReq.Body = io.NopCloser(bytes.NewReader(exec.body))
				newReq.Header = exec.ClearHopHeaders(newReq.Header)
				exec.ReplaceCustomizeTag(urlParser, bodyParser, newReq)
				result, err := exec.RunOnce(newReq)
				if isDebug {
					fmt.Printf("「DEBUG」Response Message : %s \n", result.ResponseMessage)
				}
				if err == nil {
					exec.resultChan <- result
				}
				<-concurrencyBuffer
			}()
		}
	}
OUT:
	wg.Wait()
	close(concurrencyBuffer)
	close(exec.resultChan)
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

func (exec *Executor) RunOnce(req *http.Request) (HttpTracerResult, error) {
	tracer := &HttpTracer{}
	execReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	execReq.Header = req.Header
	if err != nil {
		return HttpTracerResult{}, fmt.Errorf("执行请求失败,错误原因:%s", err.Error())
	}
	execReq = execReq.WithContext(httptrace.WithClientTrace(execReq.Context(), tracer.Trace()))
	resp, err := http.DefaultClient.Do(execReq)
	if err != nil {
		return HttpTracerResult{}, fmt.Errorf("执行请求失败,错误原因:%s", err.Error())
	}
	defer resp.Body.Close()
	return tracer.Result(execReq, resp), nil
}

func (exec *Executor) Result() *Statistic {
	return exec.statistic
}
