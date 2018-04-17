package main

import (
	"fmt"
)

func main() {
	q := make(chan int)
	fetchPara := DoParallel(func() {
		f := NewIntFetcher(q)
		f.Fetch()
	}, 2)

	go fetchPara.Wait(func() {
		close(q)
	})

	workPara := DoParallel(func() {
		w := NewIntWorker(q)
		w.Exec()
	}, 10)

	workPara.Wait(nil)
}

func NewIntFetcher(queue chan<- int) *Fetcher {
	return &Fetcher{
		Queue: queue,
	}
}

type Fetcher struct {
	Queue chan<- int
}

func (f *Fetcher) Fetch() error {
	for i := 1; i < 100; i++ {
		f.Queue <- i
	}

	return nil
}

func NewIntWorker(queue <-chan int) *Worker {
	return &Worker{
		Queue: queue,
	}
}

type Worker struct {
	Queue <-chan int
}

func (w *Worker) Exec() error {
	for i := range w.Queue {
		fmt.Printf("worker exec %d\n", i)
	}

	return nil
}
