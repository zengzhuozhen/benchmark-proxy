package core

import (
	"context"
	"sync"
)

var pool = NewPool(1)

type Pool struct {
	ctx      context.Context
	cancelFn context.CancelFunc
	cap      int
	taskChan chan Task
	taskWg   sync.WaitGroup
}

func NewPool(cap int) *Pool {
	ctx, cancelFn := context.WithCancel(context.TODO())
	p := &Pool{
		taskChan: make(chan Task),
		taskWg:   sync.WaitGroup{},
		cap:      cap,
		ctx:      ctx,
		cancelFn: cancelFn,
	}
	p.run()
	return p
}

func (p *Pool) setCap(cap int) {
	p.cap = cap
	p.cancelFn()
	ctx, cancelFn := context.WithCancel(context.TODO())
	p.ctx = ctx
	p.cancelFn = cancelFn
	p.run()
}

func (p *Pool) addTask(task Task) {
	p.taskWg.Add(1)
	p.taskChan <- task
}

func (p *Pool) run() {
	for i := 0; i < p.cap; i++ {
		go func() {
			for {
				select {
				case task := <-p.taskChan:
					task.fn()
					p.taskWg.Done()
				case _ = <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

func (p *Pool) wait() {
	p.taskWg.Wait()
}

type Task struct {
	fn func()
}

func createTask(fn func()) Task {
	return Task{fn: fn}
}
