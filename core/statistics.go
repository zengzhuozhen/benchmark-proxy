package core

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Statistic struct {
	Durations    []DurationInfo
	SuccessCount int32
	FailCount    int32
	Total        int32
	printBuf     []byte
}

func (s *Statistic) Aggregate(tracerResultChan <-chan HttpTracerResult) {
	for {
		select {
		case tracerResult, ok := <-tracerResultChan:
			if !ok {
				return
			}
			s.Durations = append(s.Durations, tracerResult.Duration)
			if tracerResult.IsSuccess {
				atomic.AddInt32(&s.SuccessCount, 1)
			} else {
				atomic.AddInt32(&s.FailCount, 1)
			}
			atomic.AddInt32(&s.Total, 1)
		}
	}
}

func (s *Statistic) Print() []byte {
	// wait for aggregate result
	time.Sleep(time.Second)
	var (
		totalDuration         int64
		totalDNSLookup        int64
		totalTCPConnection    int64
		totalTLSHandshake     int64
		totalServerProcessing int64
		totalContentTransfer  int64
	)
	for _, duration := range s.Durations {
		totalDuration += duration.Total
		totalDNSLookup += duration.DNSLookup
		totalTCPConnection += duration.TCPConnection
		totalTLSHandshake += duration.TLSHandshake
		totalServerProcessing += duration.ServerProcessing
		totalContentTransfer += duration.ContentTransfer
	}

	s.appendLine(fmt.Sprintf("Total Duration....................avg=%.2f(ms) total=%.2d(ms)", float64(totalDuration)/float64(s.Total), totalDuration))
	s.appendLine(fmt.Sprintf("----DNS Lookup....................avg=%.2f(ms)", float64(totalDNSLookup)/float64(s.Total)))
	s.appendLine(fmt.Sprintf("----TCP Connection.......... .....avg=%.2f(ms)", float64(totalTCPConnection)/float64(s.Total)))
	s.appendLine(fmt.Sprintf("----TLS Handshake.................avg=%.2f(ms)", float64(totalTLSHandshake)/float64(s.Total)))
	s.appendLine(fmt.Sprintf("----Server Processing.............avg=%.2f(ms)", float64(totalServerProcessing)/float64(s.Total)))
	s.appendLine(fmt.Sprintf("----Content Transfer..............avg=%.2f(ms)", float64(totalContentTransfer)/float64(s.Total)))

	s.appendLine(fmt.Sprintf("SuccessCount................%d", s.SuccessCount))
	s.appendLine(fmt.Sprintf("FailCount...................%d", s.FailCount))
	s.appendLine(fmt.Sprintf("Total.......................%d", s.Total))
	return s.PrintBuf()
}

func (s *Statistic) appendLine(msg string) {
	s.printBuf = append(s.printBuf, []byte(fmt.Sprint(msg+"\n"))...)
}

func (s *Statistic) PrintBuf() []byte {
	return s.printBuf
}
