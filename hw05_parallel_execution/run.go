package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n int, m int) error {
	var sharedCounter uint64
	var wgWorker sync.WaitGroup
	tasksChan := make(chan Task)

	for i := 0; i < n; i++ {
		go doWork(&wgWorker, tasksChan, m, &sharedCounter)
		wgWorker.Add(1)
	}

	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)
	wgWorker.Wait()

	if atomic.LoadUint64(&sharedCounter) > uint64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func doWork(wgWorker *sync.WaitGroup, tasksChan <-chan Task, m int, sharedCounter *uint64) {
	defer wgWorker.Done()
	for task := range tasksChan {
		if atomic.LoadUint64(sharedCounter) <= uint64(m) {
			err := task()
			if err != nil {
				atomic.AddUint64(sharedCounter, 1)
				runtime.Gosched()
			}
		}

		continue
	}
}
