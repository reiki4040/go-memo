package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", middle(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello world"))
	})))
	mux.HandleFunc("/sleep15", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(15 * time.Second)
		w.Write([]byte("wake up. I was sleeping until 15 seconds."))
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// signal handling
	signalCtx, signalCancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signalCancelFunc()

	// graceful shutdown with timeout.
	waitGraceful := make(chan struct{})
	go func() {
		<-signalCtx.Done()
		log.Print("got signal. server will do graceful shutdown...")
		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		if err := s.Shutdown(ctx); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Printf("failed graceful shutdown with timeout %v passed", time.Duration(5*time.Second))
			} else {
				log.Printf("failed graceful shutdown: %v", err)
			}
		} else {
			log.Print("OK. Server is shotdown normaly.")
		}
		close(waitGraceful)
	}()

	if err := s.ListenAndServe(); err != nil {
		log.Print(err)
	}

	/*
		got signal and call Shutdown() then ListenAndServe() returns 'Server closed'.
		it mean is stop new requests recieving, however recieved requests still
		processing with until graceful shutdown timeout.
		main() will finish if not waiting Shutdown() finished.
	*/
	<-waitGraceful
}

func middle(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		// とりあえずヘッダの有無だけ
		if auth == "" {
			http.Error(w, "Unauthorize request", http.StatusUnauthorized)
			return
		}
	}

	return http.HandlerFunc(fn)
}