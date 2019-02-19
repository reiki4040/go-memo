package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	worker := NewWorker()
	err := worker.Go()
	if err != nil {
		log.Fatalf("failed start worker: %v", err)
	}

	// wait signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sigChan

	// graceful stop or timeout
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = worker.Stop(ctx)
	if err != nil {
		log.Fatalf("failed stop: %v", err)
	} else {
		log.Printf("finished worker.")
	}
}

type Datum struct {
	Msg string
}

func NewSubscriber() *Subscriber {
	return &Subscriber{}
}

type Subscriber struct {
	quitChan chan struct{}
}

func (s *Subscriber) Start() (chan Datum, error) {
	dataChan := make(chan Datum)
	s.quitChan = make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case t := <-ticker.C:
				dataChan <- Datum{
					Msg: fmt.Sprintf("hello at %s", t),
				}
			case <-s.quitChan:
				ticker.Stop()
				close(dataChan)
				return
			}
		}
	}()

	return dataChan, nil
}

func (s *Subscriber) Stop() error {
	close(s.quitChan)
	//time.Sleep(30 * time.Second)
	return nil
}

func NewWorker() *Worker {
	return &Worker{}
}

type Worker struct {
	subscriber *Subscriber
	finishChan chan struct{}
}

func (w *Worker) Go() error {
	// finished chan notify completed processing data that got until channel closed.
	w.finishChan = make(chan struct{})

	// start subscribe
	w.subscriber = NewSubscriber()
	dataChan, err := w.subscriber.Start()
	if err != nil {
		return err
	}

	go func() {
		for d := range dataChan {
			// processing
			log.Printf("processing data: %s", d.Msg)
		}
		log.Printf("finished data processing because data channel closed.")

		// notify processed data that got until data chan closed.
		close(w.finishChan)
	}()

	return nil
}

// graceful stop or timeout
func (w *Worker) Stop(ctx context.Context) error {
	errChan := make(chan error)

	finishTask := func() error {
		err := w.subscriber.Stop()
		if err != nil {
			return err
		}

		// waiting worker graceful stop
		<-w.finishChan
		return nil
	}

	go func() {
		errChan <- finishTask()
	}()

	select {
	case <-ctx.Done():
		// got cancel
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
