package main

import (
	"sync"
)

type ParallelWG struct {
	wg *sync.WaitGroup
}

func (p *ParallelWG) Wait(allDoneCall func()) {
	p.wg.Wait()

	if allDoneCall != nil {
		allDoneCall()
	}
}

func DoParallel(fn func(), size int) *ParallelWG {
	wg := &sync.WaitGroup{}
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func() {
			fn()
			wg.Done()
		}()
	}

	return &ParallelWG{
		wg: wg,
	}
}
