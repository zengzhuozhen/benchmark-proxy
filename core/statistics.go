package core

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Statistic struct {
	TotalDuration int64
	SuccessCount  int32
	FailCount     int32
	Total         int32
	printBuf      []byte
}

func (s *Statistic) Aggregate(tracerResultChan <-chan HttpTracerResult) {
	for {
		select {
		case tracerResult, ok := <-tracerResultChan:
			if !ok {
				return
			}
			atomic.AddInt64(&s.TotalDuration, int64(tracerResult.Duration))
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
	s.appendLine(fmt.Sprintf("Duration....................avg=%.2f(ms) total=%.2f(ms)", float64(s.TotalDuration/1e6)/float64(s.Total), float64(s.TotalDuration)/1e6))
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
