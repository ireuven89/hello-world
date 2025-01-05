package test

import (
	"fmt"
	"testing"
	"time"
)

func TestServer(t *testing.T) {

	startTime := time.Now()
	// Creating a channel
	channel := make(chan int)

	// Creating 10.000 workers to execute the task
	for i := 0; i < 10000; i++ {
		go someTask(i, channel)
	}

	// Filling channel with 100.000 numbers to be executed
	for i := 0; i < 100000; i++ {
		channel <- i
	}

	t.Logf("took total: %v", time.Now().UnixMilli()-startTime.UnixMilli())
}

func someTask(id int, data chan int) {
	for taskId := range data {
		time.Sleep(2 * time.Second)
		fmt.Printf("Worker: %d executed Task %d\n", id, taskId)
	}
}
