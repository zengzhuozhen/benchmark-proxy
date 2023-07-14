package core

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {

	var i int

	task := createTask(func() {
		i++
		fmt.Println(i)
	})

	for j := 0; j < 100; j++ {
		pool.addTask(task)
	}

	pool.setCap(10)

	for j := 0; j < 200; j++ {
		pool.addTask(task)
	}

	a := make(chan struct{})
	a <- struct{}{}
}
